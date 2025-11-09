package user

import (
	"testing"
	"time"
)

func testTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestNewUser(t *testing.T) {
	// Given: 유효한 사용자 정보
	id, _ := NewUserID(1)
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("password123")

	// When: 새 사용자를 생성
	u, err := NewUser(id, email, password)

	// Then: 사용자가 성공적으로 생성됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.ID() != id {
		t.Errorf("expected ID %v, got %v", id, u.ID())
	}
	if u.Email() != email {
		t.Errorf("expected email %v, got %v", email, u.Email())
	}
	if u.Role() != UserRole() {
		t.Errorf("expected user role, got %v", u.Role())
	}
	if u.EmailVerified() {
		t.Error("expected email not verified")
	}
	if !u.Active() {
		t.Error("expected user to be active")
	}
	if len(u.DomainEvents()) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(u.DomainEvents()))
	}
}

func TestUser_Authenticate(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		active    bool
		want      bool
	}{
		{
			name:      "올바른 비밀번호와 활성 계정",
			plaintext: "password123",
			active:    true,
			want:      true,
		},
		{
			name:      "잘못된 비밀번호",
			plaintext: "wrongpassword",
			active:    true,
			want:      false,
		},
		{
			name:      "비활성화된 계정",
			plaintext: "password123",
			active:    false,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 사용자
			id, _ := NewUserID(1)
			email, _ := NewEmail("test@example.com")
			password, _ := NewPassword("password123")
			u, _ := NewUser(id, email, password)

			if !tt.active {
				u.Deactivate()
			}

			// When: 인증 시도
			got := u.Authenticate(tt.plaintext)

			// Then: 예상된 결과
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestUser_VerifyEmail(t *testing.T) {
	t.Run("이메일 인증 성공", func(t *testing.T) {
		// Given: 인증되지 않은 사용자
		id, _ := NewUserID(1)
		email, _ := NewEmail("test@example.com")
		password, _ := NewPassword("password123")
		u, _ := NewUser(id, email, password)
		u.ClearEvents()

		// When: 이메일 인증
		err := u.VerifyEmail()

		// Then: 이메일이 인증됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !u.EmailVerified() {
			t.Error("expected email to be verified")
		}
		if len(u.DomainEvents()) != 1 {
			t.Errorf("expected 1 domain event, got %d", len(u.DomainEvents()))
		}
	})

	t.Run("이미 인증된 이메일", func(t *testing.T) {
		// Given: 이미 인증된 사용자
		id, _ := NewUserID(1)
		email, _ := NewEmail("test@example.com")
		password, _ := NewPassword("password123")
		u, _ := NewUser(id, email, password)
		u.VerifyEmail()
		u.ClearEvents()

		// When: 다시 이메일 인증
		err := u.VerifyEmail()

		// Then: 에러 없이 처리됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(u.DomainEvents()) != 0 {
			t.Errorf("expected no domain events, got %d", len(u.DomainEvents()))
		}
	})
}

func TestUser_ChangePassword(t *testing.T) {
	// Given: 사용자
	id, _ := NewUserID(1)
	email, _ := NewEmail("test@example.com")
	oldPassword, _ := NewPassword("oldpassword123")
	u, _ := NewUser(id, email, oldPassword)
	u.ClearEvents()

	newPassword, _ := NewPassword("newpassword123")

	// When: 비밀번호 변경
	err := u.ChangePassword(newPassword)

	// Then: 비밀번호가 변경됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Password() != newPassword {
		t.Error("expected password to be updated")
	}
	if !u.Authenticate("newpassword123") {
		t.Error("expected to authenticate with new password")
	}
	if u.Authenticate("oldpassword123") {
		t.Error("expected not to authenticate with old password")
	}
	if len(u.DomainEvents()) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(u.DomainEvents()))
	}
}

func TestUser_ChangeRole(t *testing.T) {
	t.Run("사용자에서 관리자로 변경", func(t *testing.T) {
		// Given: 사용자 역할을 가진 사용자
		id, _ := NewUserID(1)
		email, _ := NewEmail("test@example.com")
		password, _ := NewPassword("password123")
		u, _ := NewUser(id, email, password)
		u.ClearEvents()

		// When: 관리자로 역할 변경
		err := u.ChangeRole(AdminRole())

		// Then: 역할이 변경됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if u.Role() != AdminRole() {
			t.Errorf("expected role %v, got %v", AdminRole(), u.Role())
		}
		if len(u.DomainEvents()) != 1 {
			t.Errorf("expected 1 domain event, got %d", len(u.DomainEvents()))
		}
	})

	t.Run("관리자에서 사용자로 변경", func(t *testing.T) {
		// Given: 관리자 역할을 가진 사용자
		id, _ := NewUserID(1)
		email, _ := NewEmail("test@example.com")
		password, _ := NewPassword("password123")
		u := ReconstructUser(id, email, password, AdminRole(), false, true, testTime(), testTime())

		// When: 사용자로 역할 변경
		err := u.ChangeRole(UserRole())

		// Then: 역할이 변경됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if u.Role() != UserRole() {
			t.Errorf("expected role %v, got %v", UserRole(), u.Role())
		}
		if len(u.DomainEvents()) != 1 {
			t.Errorf("expected 1 domain event, got %d", len(u.DomainEvents()))
		}
	})

	t.Run("동일한 역할로 변경", func(t *testing.T) {
		// Given: 사용자
		id, _ := NewUserID(1)
		email, _ := NewEmail("test@example.com")
		password, _ := NewPassword("password123")
		u, _ := NewUser(id, email, password)
		u.ClearEvents()

		// When: 동일한 역할로 변경
		err := u.ChangeRole(UserRole())

		// Then: 변경 없음
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(u.DomainEvents()) != 0 {
			t.Errorf("expected no domain events, got %d", len(u.DomainEvents()))
		}
	})
}

func TestUser_Deactivate(t *testing.T) {
	// Given: 활성 사용자
	id, _ := NewUserID(1)
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("password123")
	u, _ := NewUser(id, email, password)
	u.ClearEvents()

	// When: 비활성화
	err := u.Deactivate()

	// Then: 비활성화됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Active() {
		t.Error("expected user to be inactive")
	}
	if len(u.DomainEvents()) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(u.DomainEvents()))
	}
}

func TestUser_Activate(t *testing.T) {
	// Given: 비활성 사용자
	id, _ := NewUserID(1)
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("password123")
	u, _ := NewUser(id, email, password)
	u.Deactivate()
	u.ClearEvents()

	// When: 활성화
	err := u.Activate()

	// Then: 활성화됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !u.Active() {
		t.Error("expected user to be active")
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{
			name: "관리자 역할",
			role: AdminRole(),
			want: true,
		},
		{
			name: "사용자 역할",
			role: UserRole(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 특정 역할을 가진 사용자
			id, _ := NewUserID(1)
			email, _ := NewEmail("test@example.com")
			password, _ := NewPassword("password123")
			u := ReconstructUser(id, email, password, tt.role, false, true, testTime(), testTime())

			// When: 관리자 확인
			got := u.IsAdmin()

			// Then: 예상된 결과
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
