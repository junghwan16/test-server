package user

import (
	"github.com/junghwan16/test-server/internal/shared/domain"
)

// UserRegistered is fired when a new user registers
type UserRegistered struct {
	domain.BaseEvent
	UserID UserID
	Email  Email
}

// EventType returns the event type
func (e UserRegistered) EventType() string {
	return "identity.user.registered"
}

// NewUserRegistered creates a new UserRegistered event
func NewUserRegistered(userID UserID, email Email) UserRegistered {
	return UserRegistered{
		BaseEvent: domain.NewBaseEvent(),
		UserID:    userID,
		Email:     email,
	}
}

// EmailVerified is fired when an email is verified
type EmailVerified struct {
	domain.BaseEvent
	UserID UserID
}

// EventType returns the event type
func (e EmailVerified) EventType() string {
	return "identity.user.email_verified"
}

// NewEmailVerified creates a new EmailVerified event
func NewEmailVerified(userID UserID) EmailVerified {
	return EmailVerified{
		BaseEvent: domain.NewBaseEvent(),
		UserID:    userID,
	}
}

// PasswordChanged is fired when a password is changed
type PasswordChanged struct {
	domain.BaseEvent
	UserID UserID
}

// EventType returns the event type
func (e PasswordChanged) EventType() string {
	return "identity.user.password_changed"
}

// NewPasswordChanged creates a new PasswordChanged event
func NewPasswordChanged(userID UserID) PasswordChanged {
	return PasswordChanged{
		BaseEvent: domain.NewBaseEvent(),
		UserID:    userID,
	}
}

// UserDeactivated is fired when a user is deactivated
type UserDeactivated struct {
	domain.BaseEvent
	UserID UserID
}

// EventType returns the event type
func (e UserDeactivated) EventType() string {
	return "identity.user.deactivated"
}

// NewUserDeactivated creates a new UserDeactivated event
func NewUserDeactivated(userID UserID) UserDeactivated {
	return UserDeactivated{
		BaseEvent: domain.NewBaseEvent(),
		UserID:    userID,
	}
}

// RoleChanged is fired when a user's role is changed
type RoleChanged struct {
	domain.BaseEvent
	UserID  UserID
	OldRole Role
	NewRole Role
}

// EventType returns the event type
func (e RoleChanged) EventType() string {
	return "identity.user.role_changed"
}

// NewRoleChanged creates a new RoleChanged event
func NewRoleChanged(userID UserID, oldRole, newRole Role) RoleChanged {
	return RoleChanged{
		BaseEvent: domain.NewBaseEvent(),
		UserID:    userID,
		OldRole:   oldRole,
		NewRole:   newRole,
	}
}
