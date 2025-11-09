package handler

import (
	"context"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

type contextKey string

const (
	contextKeyUser contextKey = "user"
)

// SetUserInContext sets the user in the context
func SetUserInContext(ctx context.Context, u *user.User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

// GetUserFromContext retrieves the user from context
func GetUserFromContext(ctx context.Context) *user.User {
	u, ok := ctx.Value(contextKeyUser).(*user.User)
	if !ok {
		return nil
	}
	return u
}
