package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/junghwan16/test-server/internal/identity/domain/session"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// RedisSessionRepository implements session.Repository using Redis
type RedisSessionRepository struct {
	client *redis.Client
}

// NewRedisSessionRepository creates a new Redis-based SessionRepository
func NewRedisSessionRepository(client *redis.Client) *RedisSessionRepository {
	return &RedisSessionRepository{client: client}
}

type redisSessionData struct {
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RedisSessionRepository) Save(s *session.Session) error {
	ctx := context.Background()

	data := redisSessionData{
		UserID:    s.UserID().Value(),
		ExpiresAt: s.ExpiresAt(),
		CreatedAt: s.CreatedAt(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := sessionKey(s.ID().Value())
	ttl := time.Until(s.ExpiresAt())
	if ttl <= 0 {
		return errors.New("session already expired")
	}

	return r.client.Set(ctx, key, jsonData, ttl).Err()
}

func (r *RedisSessionRepository) FindByID(id session.SessionID) (*session.Session, error) {
	ctx := context.Background()

	key := sessionKey(id.Value())
	jsonData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("session not found or expired")
		}
		return nil, err
	}

	var data redisSessionData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	// Check if expired (double check, though Redis should auto-expire)
	if time.Now().After(data.ExpiresAt) {
		r.client.Del(ctx, key) // cleanup
		return nil, errors.New("session not found or expired")
	}

	userID, _ := user.NewUserID(data.UserID)
	return session.ReconstructSession(id, userID, data.ExpiresAt, data.CreatedAt), nil
}

func (r *RedisSessionRepository) Delete(id session.SessionID) error {
	ctx := context.Background()
	key := sessionKey(id.Value())
	return r.client.Del(ctx, key).Err()
}

// DeleteExpired is a no-op for Redis since it handles expiration automatically
func (r *RedisSessionRepository) DeleteExpired() error {
	// Redis automatically removes expired keys
	return nil
}

// DeleteByUserID removes all sessions for a given user
func (r *RedisSessionRepository) DeleteByUserID(userID user.UserID) error {
	ctx := context.Background()

	// Use a scan to find all session keys
	pattern := "session:*"
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		jsonData, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var data redisSessionData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			continue
		}

		if data.UserID == userID.Value() {
			r.client.Del(ctx, key)
		}
	}

	return iter.Err()
}

func sessionKey(id string) string {
	return "session:" + id
}
