package http

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/junghwan16/test-server/internal/auth/domain"
)

type JWTEncoder struct {
	secretKey string
}

func NewJWTEncoder(secretKey string) *JWTEncoder {
	return &JWTEncoder{
		secretKey: secretKey,
	}
}

func (e *JWTEncoder) EncodeSession(session *domain.Session) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   session.UserID().String(),
		ExpiresAt: jwt.NewNumericDate(session.ExpiresAt()),
		IssuedAt:  jwt.NewNumericDate(session.CreatedAt()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(e.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

func (e *JWTEncoder) DecodeToken(tokenString string) (userID string, err error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(e.secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid JWT: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid JWT")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("JWT expired")
	}

	return claims.Subject, nil
}
