package domain

import "errors"

var (
	// ErrUserNotFound는 사용자를 찾을 수 없을 때 반환되는 에러입니다.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials는 인증 정보가 올바르지 않을 때 반환되는 에러입니다.
	// 보안상 이유로 사용자명이 존재하지 않는지, 비밀번호가 틀렸는지 구분하지 않습니다.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUsernameAlreadyExists는 이미 존재하는 사용자명으로 등록을 시도할 때 반환되는 에러입니다.
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// ErrUsernameTaken는 사용자명이 이미 사용 중일 때 반환되는 에러입니다.
	ErrUsernameTaken = errors.New("username already taken")

	// ErrInvalidPassword는 비밀번호가 올바르지 않을 때 반환되는 에러입니다.
	ErrInvalidPassword = errors.New("invalid password")
)
