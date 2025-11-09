package user

import (
	"time"

	"github.com/junghwan16/test-server/internal/shared/domain"
)

// User is the aggregate root for user identity
type User struct {
	id            UserID
	email         Email
	password      Password
	role          Role
	emailVerified bool
	active        bool
	createdAt     time.Time
	updatedAt     time.Time

	// Domain events
	events []domain.DomainEvent
}

// NewUser creates a new User aggregate (factory method)
func NewUser(id UserID, email Email, password Password) (*User, error) {
	user := &User{
		id:            id,
		email:         email,
		password:      password,
		role:          UserRole(),
		emailVerified: false,
		active:        true,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
		events:        make([]domain.DomainEvent, 0),
	}

	user.addEvent(NewUserRegistered(id, email))

	return user, nil
}

// ReconstructUser reconstructs a User from persistence (not a new registration)
func ReconstructUser(
	id UserID,
	email Email,
	password Password,
	role Role,
	emailVerified, active bool,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:            id,
		email:         email,
		password:      password,
		role:          role,
		emailVerified: emailVerified,
		active:        active,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		events:        make([]domain.DomainEvent, 0),
	}
}

func (u *User) ID() UserID           { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) Password() Password   { return u.password }
func (u *User) Role() Role           { return u.role }
func (u *User) EmailVerified() bool  { return u.emailVerified }
func (u *User) Active() bool         { return u.active }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Authenticate checks if the password is correct
func (u *User) Authenticate(plaintext string) bool {
	if !u.active {
		return false
	}
	return u.password.Matches(plaintext)
}

// VerifyEmail marks the email as verified
func (u *User) VerifyEmail() error {
	if u.emailVerified {
		return nil // Already verified
	}

	u.emailVerified = true
	u.updatedAt = time.Now()
	u.addEvent(NewEmailVerified(u.id))

	return nil
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(newPassword Password) error {
	u.password = newPassword
	u.updatedAt = time.Now()
	u.addEvent(NewPasswordChanged(u.id))

	return nil
}

// ChangeRole changes the user's role (admin operation)
func (u *User) ChangeRole(newRole Role) error {
	if u.role.Equals(newRole) {
		return nil // No change
	}

	oldRole := u.role
	u.role = newRole
	u.updatedAt = time.Now()
	u.addEvent(NewRoleChanged(u.id, oldRole, newRole))

	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() error {
	if !u.active {
		return nil // Already deactivated
	}

	u.active = false
	u.updatedAt = time.Now()
	u.addEvent(NewUserDeactivated(u.id))

	return nil
}

// Activate activates the user account
func (u *User) Activate() error {
	u.active = true
	u.updatedAt = time.Now()
	return nil
}

// IsAdmin returns true if the user is an admin
func (u *User) IsAdmin() bool {
	return u.role.IsAdmin()
}

func (u *User) addEvent(event domain.DomainEvent) {
	u.events = append(u.events, event)
}

// DomainEvents returns all uncommitted domain events
func (u *User) DomainEvents() []domain.DomainEvent {
	return u.events
}

// ClearEvents clears all domain events (after publishing)
func (u *User) ClearEvents() {
	u.events = make([]domain.DomainEvent, 0)
}
