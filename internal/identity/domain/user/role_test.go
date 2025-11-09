package user

import "testing"

func TestNewRole(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "사용자 역할",
			value:   "user",
			wantErr: false,
		},
		{
			name:    "관리자 역할",
			value:   "admin",
			wantErr: false,
		},
		{
			name:    "대문자 역할은 에러",
			value:   "ADMIN",
			wantErr: true,
		},
		{
			name:    "잘못된 역할",
			value:   "invalid",
			wantErr: true,
		},
		{
			name:    "빈 역할",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: 역할 생성
			role, err := NewRole(tt.value)

			// Then: 예상된 결과
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				if role.Value() == "" {
					t.Error("expected non-empty role value")
				}
			}
		})
	}
}

func TestRole_IsAdmin(t *testing.T) {
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
			// When: 관리자 확인
			got := tt.role.IsAdmin()

			// Then: 예상된 결과
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestRole_IsUser(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{
			name: "사용자 역할",
			role: UserRole(),
			want: true,
		},
		{
			name: "관리자 역할",
			role: AdminRole(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: 사용자 확인
			got := tt.role.IsUser()

			// Then: 예상된 결과
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestRole_Equals(t *testing.T) {
	// Given: 두 역할
	role1 := UserRole()
	role2 := UserRole()
	role3 := AdminRole()

	// When & Then: 같은 역할 비교
	if !role1.Equals(role2) {
		t.Error("expected equal roles")
	}

	// When & Then: 다른 역할 비교
	if role1.Equals(role3) {
		t.Error("expected not equal roles")
	}
}
