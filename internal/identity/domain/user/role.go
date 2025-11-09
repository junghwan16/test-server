package user

import (
	"errors"
)

// Role is a value object representing a user's role
type Role struct {
	value string
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

var ErrInvalidRole = errors.New("invalid role")

// NewRole creates a new Role
func NewRole(value string) (Role, error) {
	if value != RoleUser && value != RoleAdmin {
		return Role{}, ErrInvalidRole
	}
	return Role{value: value}, nil
}

// UserRole returns the user role
func UserRole() Role {
	return Role{value: RoleUser}
}

// AdminRole returns the admin role
func AdminRole() Role {
	return Role{value: RoleAdmin}
}

// Value returns the role string
func (r Role) Value() string {
	return r.value
}

// IsAdmin returns true if the role is admin
func (r Role) IsAdmin() bool {
	return r.value == RoleAdmin
}

// IsUser returns true if the role is user
func (r Role) IsUser() bool {
	return r.value == RoleUser
}

// Equals checks if two roles are equal
func (r Role) Equals(other Role) bool {
	return r.value == other.value
}

// String returns string representation
func (r Role) String() string {
	return r.value
}
