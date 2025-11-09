package user

// Repository defines the interface for User aggregate persistence
type Repository interface {
	// NextID generates a new UserID
	NextID() UserID

	Save(user *User) error

	// FindByID retrieves a User by ID
	FindByID(id UserID) (*User, error)

	// FindByEmail retrieves a User by email
	FindByEmail(email Email) (*User, error)

	// FindAll retrieves all users with pagination
	FindAll(limit, offset int) ([]*User, int64, error)

	Delete(id UserID) error
}
