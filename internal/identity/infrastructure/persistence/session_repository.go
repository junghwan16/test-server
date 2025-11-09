package persistence

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/identity/domain/session"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// SessionModel is the GORM model for Session aggregate
type SessionModel struct {
	ID        string    `gorm:"primarykey"`
	UserID    uint      `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

// TableName specifies the table name
func (SessionModel) TableName() string {
	return "sessions"
}

// SessionRepository implements session.Repository using GORM
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Save(s *session.Session) error {
	model := r.toModel(s)
	return r.db.Create(&model).Error
}

// FindByID retrieves a Session by ID
func (r *SessionRepository) FindByID(id session.SessionID) (*session.Session, error) {
	var model SessionModel
	err := r.db.Where("id = ? AND expires_at > ?", id.Value(), time.Now()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found or expired")
		}
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *SessionRepository) Delete(id session.SessionID) error {
	return r.db.Delete(&SessionModel{}, "id = ?", id.Value()).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&SessionModel{}).Error
}

func (r *SessionRepository) DeleteByUserID(userID user.UserID) error {
	return r.db.Where("user_id = ?", userID.Value()).Delete(&SessionModel{}).Error
}

// Mapping

func (r *SessionRepository) toModel(s *session.Session) SessionModel {
	return SessionModel{
		ID:        s.ID().Value(),
		UserID:    s.UserID().Value(),
		ExpiresAt: s.ExpiresAt(),
		CreatedAt: s.CreatedAt(),
	}
}

func (r *SessionRepository) toDomain(m *SessionModel) *session.Session {
	id, _ := session.NewSessionID(m.ID)
	userID, _ := user.NewUserID(m.UserID)

	return session.ReconstructSession(id, userID, m.ExpiresAt, m.CreatedAt)
}
