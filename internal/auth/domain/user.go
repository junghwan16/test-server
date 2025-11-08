package domain

import "errors"

// User는 시스템의 사용자 Aggregate Root입니다.
// 불변성을 보장하기 위해 필드를 private으로 선언합니다.
type User struct {
	id       UserID
	username Username
	password HashedPassword
}

// NewUser는 새로운 사용자를 생성합니다.
// 이 함수는 회원가입 시 사용됩니다.
func NewUser(id UserID, username Username, plainPassword string) (*User, error) {
	hashedPassword, err := NewHashedPasswordFromPlain(plainPassword)
	if err != nil {
		return nil, err
	}

	return &User{
		id:       id,
		username: username,
		password: hashedPassword,
	}, nil
}

// ReconstructUser는 저장소에서 읽어온 데이터로부터 User를 재구성합니다.
// 이미 해시된 비밀번호를 사용합니다.
func ReconstructUser(id UserID, username Username, passwordHash string) (*User, error) {
	hashedPassword, err := NewHashedPasswordFromHash(passwordHash)
	if err != nil {
		return nil, err
	}

	return &User{
		id:       id,
		username: username,
		password: hashedPassword,
	}, nil
}

// ID는 사용자의 식별자를 반환합니다.
func (u *User) ID() UserID {
	return u.id
}

// Username은 사용자명을 반환합니다.
func (u *User) Username() Username {
	return u.username
}

// PasswordHash는 저장을 위한 비밀번호 해시를 반환합니다.
func (u *User) PasswordHash() string {
	return u.password.Hash()
}

// Authenticate는 사용자를 인증합니다.
// 도메인 로직: 사용자는 자신의 비밀번호를 검증할 수 있어야 합니다.
func (u *User) Authenticate(plainPassword string) error {
	if err := u.password.Verify(plainPassword); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

// VerifyPassword는 비밀번호가 올바른지 확인합니다.
func (u *User) VerifyPassword(plainPassword string) bool {
	return u.password.Verify(plainPassword) == nil
}

// ChangeUsername은 사용자명을 변경합니다.
func (u *User) ChangeUsername(newUsername Username) {
	u.username = newUsername
}

// ChangePassword는 비밀번호를 변경합니다.
func (u *User) ChangePassword(plainPassword string) error {
	hashedPassword, err := NewHashedPasswordFromPlain(plainPassword)
	if err != nil {
		return err
	}
	u.password = hashedPassword
	return nil
}
