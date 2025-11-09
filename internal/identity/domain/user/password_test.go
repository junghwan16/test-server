package user

import "testing"

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		wantErr   bool
	}{
		{
			name:      "유효한 비밀번호",
			plaintext: "password123",
			wantErr:   false,
		},
		{
			name:      "8자 정확히",
			plaintext: "12345678",
			wantErr:   false,
		},
		{
			name:      "짧은 비밀번호",
			plaintext: "short",
			wantErr:   true,
		},
		{
			name:      "빈 비밀번호",
			plaintext: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: 비밀번호 생성
			password, err := NewPassword(tt.plaintext)

			// Then: 예상된 결과
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				if password.Hash() == "" {
					t.Error("expected non-empty hash")
				}
				if password.Hash() == tt.plaintext {
					t.Error("expected hash to be different from plaintext")
				}
			}
		})
	}
}

func TestPassword_Matches(t *testing.T) {
	tests := []struct {
		name      string
		original  string
		plaintext string
		want      bool
	}{
		{
			name:      "일치하는 비밀번호",
			original:  "password123",
			plaintext: "password123",
			want:      true,
		},
		{
			name:      "일치하지 않는 비밀번호",
			original:  "password123",
			plaintext: "wrongpassword",
			want:      false,
		},
		{
			name:      "빈 비밀번호",
			original:  "password123",
			plaintext: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 비밀번호
			password, _ := NewPassword(tt.original)

			// When: 비밀번호 일치 확인
			got := password.Matches(tt.plaintext)

			// Then: 예상된 결과
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestPassword_Change(t *testing.T) {
	// Given: 기존 비밀번호
	oldPassword, _ := NewPassword("oldpassword123")

	// When: 새 비밀번호로 변경
	newPassword, err := oldPassword.Change("newpassword123")

	// Then: 변경 성공
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if newPassword.Hash() == oldPassword.Hash() {
		t.Error("expected different hash")
	}
	if !newPassword.Matches("newpassword123") {
		t.Error("expected new password to match")
	}
	if newPassword.Matches("oldpassword123") {
		t.Error("expected new password not to match old password")
	}
}
