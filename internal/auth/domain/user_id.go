package domain

import "strconv"

type UserID struct {
	value uint
}

func NewUserID(value uint) UserID {
	return UserID{value: value}
}

func NewUserIDFromString(s string) (UserID, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return UserID{}, err
	}
	return UserID{value: uint(id)}, nil
}

func (u UserID) Value() uint {
	return u.value
}

func (u UserID) String() string {
	return strconv.FormatUint(uint64(u.value), 10)
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

func (u UserID) IsZero() bool {
	return u.value == 0
}
