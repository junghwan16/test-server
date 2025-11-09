package application

import (
	"errors"
	"testing"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

type mockUserRepository struct {
	users   map[uint]*user.User
	nextID  uint
	findErr error
	saveErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[uint]*user.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Save(u *user.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.users[u.ID().Value()] = u
	return nil
}

func (m *mockUserRepository) FindByID(id user.UserID) (*user.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	u, ok := m.users[id.Value()]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepository) FindByEmail(email user.Email) (*user.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, u := range m.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) FindAll(limit, offset int) ([]*user.User, int64, error) {
	if m.findErr != nil {
		return nil, 0, m.findErr
	}
	var result []*user.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, int64(len(result)), nil
}

func (m *mockUserRepository) Delete(id user.UserID) error {
	delete(m.users, id.Value())
	return nil
}

func (m *mockUserRepository) NextID() user.UserID {
	id, _ := user.NewUserID(m.nextID)
	m.nextID++
	return id
}

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("성공적으로 사용자 등록", func(t *testing.T) {
		// Given: 유효한 이메일과 비밀번호
		repo := newMockUserRepository()
		svc := NewUserService(repo)

		email := "test@example.com"
		password := "password123"

		// When: 사용자 등록
		u, err := svc.RegisterUser(email, password)

		// Then: 사용자가 생성됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if u == nil {
			t.Fatal("expected user, got nil")
		}
		if u.Email().Value() != email {
			t.Errorf("expected email %s, got %s", email, u.Email().Value())
		}
		if len(repo.users) != 1 {
			t.Errorf("expected 1 user in repo, got %d", len(repo.users))
		}
	})

	t.Run("이미 존재하는 이메일", func(t *testing.T) {
		// Given: 이미 등록된 이메일
		repo := newMockUserRepository()
		svc := NewUserService(repo)

		email := "test@example.com"
		svc.RegisterUser(email, "password123")

		// When: 같은 이메일로 재등록
		_, err := svc.RegisterUser(email, "password456")

		// Then: 에러 발생
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("잘못된 이메일 형식", func(t *testing.T) {
		// Given: 잘못된 이메일
		repo := newMockUserRepository()
		svc := NewUserService(repo)

		// When: 잘못된 이메일로 등록
		_, err := svc.RegisterUser("invalid-email", "password123")

		// Then: 에러 발생
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("짧은 비밀번호", func(t *testing.T) {
		// Given: 짧은 비밀번호
		repo := newMockUserRepository()
		svc := NewUserService(repo)

		// When: 짧은 비밀번호로 등록
		_, err := svc.RegisterUser("test@example.com", "short")

		// Then: 에러 발생
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("사용자 조회 성공", func(t *testing.T) {
		// Given: 등록된 사용자
		repo := newMockUserRepository()
		svc := NewUserService(repo)
		created, _ := svc.RegisterUser("test@example.com", "password123")

		// When: 사용자 조회
		u, err := svc.GetUser(created.ID().Value())

		// Then: 사용자를 찾음
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if u.ID() != created.ID() {
			t.Errorf("expected ID %v, got %v", created.ID(), u.ID())
		}
	})

	t.Run("존재하지 않는 사용자", func(t *testing.T) {
		// Given: 빈 저장소
		repo := newMockUserRepository()
		svc := NewUserService(repo)

		// When: 존재하지 않는 사용자 조회
		_, err := svc.GetUser(999)

		// Then: 에러 발생
		if err != ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserService_ChangePassword(t *testing.T) {
	// Given: 등록된 사용자
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	u, _ := svc.RegisterUser("test@example.com", "oldpassword123")

	// When: 비밀번호 변경
	err := svc.ChangePassword(u.ID().Value(), "newpassword123")

	// Then: 비밀번호가 변경됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := svc.GetUser(u.ID().Value())
	if !updated.Authenticate("newpassword123") {
		t.Error("expected to authenticate with new password")
	}
	if updated.Authenticate("oldpassword123") {
		t.Error("expected not to authenticate with old password")
	}
}

func TestUserService_VerifyEmail(t *testing.T) {
	// Given: 등록된 사용자
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	u, _ := svc.RegisterUser("test@example.com", "password123")

	// When: 이메일 인증
	err := svc.VerifyEmail(u.ID().Value())

	// Then: 이메일이 인증됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := svc.GetUser(u.ID().Value())
	if !updated.EmailVerified() {
		t.Error("expected email to be verified")
	}
}

func TestUserService_ChangeRole(t *testing.T) {
	// Given: 등록된 사용자
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	u, _ := svc.RegisterUser("test@example.com", "password123")

	// When: 역할을 관리자로 변경
	err := svc.ChangeRole(u.ID().Value(), "admin")

	// Then: 역할이 변경됨
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := svc.GetUser(u.ID().Value())
	if !updated.IsAdmin() {
		t.Error("expected user to be admin")
	}
}

func TestUserService_SetActive(t *testing.T) {
	t.Run("사용자 비활성화", func(t *testing.T) {
		// Given: 등록된 사용자
		repo := newMockUserRepository()
		svc := NewUserService(repo)
		u, _ := svc.RegisterUser("test@example.com", "password123")

		// When: 비활성화
		err := svc.SetActive(u.ID().Value(), false)

		// Then: 비활성화됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updated, _ := svc.GetUser(u.ID().Value())
		if updated.Active() {
			t.Error("expected user to be inactive")
		}
	})

	t.Run("사용자 활성화", func(t *testing.T) {
		// Given: 비활성 사용자
		repo := newMockUserRepository()
		svc := NewUserService(repo)
		u, _ := svc.RegisterUser("test@example.com", "password123")
		svc.SetActive(u.ID().Value(), false)

		// When: 활성화
		err := svc.SetActive(u.ID().Value(), true)

		// Then: 활성화됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updated, _ := svc.GetUser(u.ID().Value())
		if !updated.Active() {
			t.Error("expected user to be active")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("다른 사용자 삭제", func(t *testing.T) {
		// Given: 두 명의 사용자
		repo := newMockUserRepository()
		svc := NewUserService(repo)
		u1, _ := svc.RegisterUser("user1@example.com", "password123")
		u2, _ := svc.RegisterUser("user2@example.com", "password123")

		// When: 다른 사용자 삭제
		err := svc.DeleteUser(u2.ID().Value(), u1.ID().Value())

		// Then: 삭제됨
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = svc.GetUser(u2.ID().Value())
		if err != ErrUserNotFound {
			t.Error("expected user to be deleted")
		}
	})

	t.Run("자기 자신 삭제 시도", func(t *testing.T) {
		// Given: 사용자
		repo := newMockUserRepository()
		svc := NewUserService(repo)
		u, _ := svc.RegisterUser("test@example.com", "password123")

		// When: 자기 자신 삭제 시도
		err := svc.DeleteUser(u.ID().Value(), u.ID().Value())

		// Then: 에러 발생
		if err != ErrCannotDeleteSelf {
			t.Errorf("expected ErrCannotDeleteSelf, got %v", err)
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	// Given: 여러 사용자
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	svc.RegisterUser("user1@example.com", "password123")
	svc.RegisterUser("user2@example.com", "password123")
	svc.RegisterUser("user3@example.com", "password123")

	// When: 사용자 목록 조회
	users, total, err := svc.ListUsers(10, 0)

	// Then: 모든 사용자 반환
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}
