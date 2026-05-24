package domain

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name:    "valid user",
			user:    User{Username: "admin", Email: "admin@mail.com"},
			wantErr: nil,
		},
		{
			name:    "empty email",
			user:    User{Username: "admin", Email: ""},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "whitespace email",
			user:    User{Username: "admin", Email: "  "},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "username too short",
			user:    User{Username: "ad", Email: "admin@mail.com"},
			wantErr: ErrInvalidName,
		},
		{
			name:    "whitespace username",
			user:    User{Username: "  ", Email: "admin@mail.com"},
			wantErr: ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
