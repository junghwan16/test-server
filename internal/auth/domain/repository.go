package domain

type UserRepository interface {
	FindByID(id UserID) (*User, error)
	FindByUsername(username Username) (*User, error)
	Save(user *User) error
	Update(user *User) error
	Delete(id UserID) error
	ExistsByUsername(username Username) (bool, error)
}
