package user

import "testing"

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "유효한 이메일",
			value:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "대문자 이메일은 소문자로 변환",
			value:   "Test@Example.COM",
			wantErr: false,
		},
		{
			name:    "공백이 있는 이메일은 trim됨",
			value:   "  test@example.com  ",
			wantErr: false,
		},
		{
			name:    "@가 없는 이메일",
			value:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "빈 문자열",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: 이메일 생성
			email, err := NewEmail(tt.value)

			// Then: 예상된 결과
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				if email.Value() == "" {
					t.Error("expected non-empty email value")
				}
			}
		})
	}
}

func TestEmail_Value(t *testing.T) {
	// Given: 유효한 이메일
	input := "Test@Example.COM"
	email, _ := NewEmail(input)

	// When: 값 조회
	got := email.Value()

	// Then: 소문자로 변환됨
	want := "test@example.com"
	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}
