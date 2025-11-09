package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/junghwan16/test-server/internal/identity/application"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

type UsersHandler struct {
	userSvc *application.UserService
}

func NewUsersHandler(userSvc *application.UserService) *UsersHandler {
	return &UsersHandler{
		userSvc: userSvc,
	}
}

// ListUsers returns a paginated list of users (admin only)
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, total, err := h.userSvc.ListUsers(limit, offset)
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	// Convert to DTOs
	userDTOs := make([]map[string]any, len(users))
	for i, u := range users {
		userDTOs[i] = userToDTO(u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"users":  userDTOs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUser returns a single user (admin only)
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	u, err := h.userSvc.GetUser(uint(id))
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userToDTO(u))
}

// UpdateUser updates a user (admin only)
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Role          *string `json:"role,omitempty"`
		Active        *bool   `json:"active,omitempty"`
		EmailVerified *bool   `json:"email_verified,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update role if provided
	if req.Role != nil {
		if err := h.userSvc.ChangeRole(uint(id), *req.Role); err != nil {
			if errors.Is(err, user.ErrInvalidRole) {
				http.Error(w, "Invalid role", http.StatusBadRequest)
			} else {
				http.Error(w, "Failed to update role", http.StatusInternalServerError)
			}
			return
		}
	}

	// Update active status if provided
	if req.Active != nil {
		if err := h.userSvc.SetActive(uint(id), *req.Active); err != nil {
			http.Error(w, "Failed to update active status", http.StatusInternalServerError)
			return
		}
	}

	// Update email verified if provided
	if req.EmailVerified != nil && *req.EmailVerified {
		if err := h.userSvc.VerifyEmail(uint(id)); err != nil {
			http.Error(w, "Failed to verify email", http.StatusInternalServerError)
			return
		}
	}

	// Get updated user
	u, err := h.userSvc.GetUser(uint(id))
	if err != nil {
		http.Error(w, "Failed to get updated user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user": userToDTO(u),
	})
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	currentUser := GetUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.userSvc.DeleteUser(uint(id), currentUser.ID().Value()); err != nil {
		if errors.Is(err, application.ErrCannotDeleteSelf) {
			http.Error(w, "Cannot delete yourself", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted",
	})
}
