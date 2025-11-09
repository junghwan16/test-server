package handler

import "github.com/junghwan16/test-server/internal/identity/domain/user"

// userToDTO converts User aggregate to DTO
func userToDTO(u *user.User) map[string]any {
	return map[string]any{
		"id":             u.ID().Value(),
		"email":          u.Email().Value(),
		"role":           u.Role().Value(),
		"email_verified": u.EmailVerified(),
		"active":         u.Active(),
		"created_at":     u.CreatedAt(),
		"updated_at":     u.UpdatedAt(),
	}
}
