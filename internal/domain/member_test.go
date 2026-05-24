package domain

import (
	"testing"
)

func TestMember_Validate(t *testing.T) {
	tests := []struct {
		name    string
		member  Member
		wantErr error
	}{
		{
			name:    "valid member",
			member:  Member{Name: "Alex", Email: "alex@mail.com", Plan: "Gold"},
			wantErr: nil,
		},
		{
			name:    "empty email",
			member:  Member{Name: "Alex", Email: "", Plan: "Gold"},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "whitespace email",
			member:  Member{Name: "Alex", Email: "  ", Plan: "Gold"},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "name too short",
			member:  Member{Name: "Al", Email: "alex@mail.com", Plan: "Gold"},
			wantErr: ErrInvalidName,
		},
		{
			name:    "whitespace name",
			member:  Member{Name: "  ", Email: "alex@mail.com", Plan: "Gold"},
			wantErr: ErrInvalidName,
		},
		{
			name:    "empty plan",
			member:  Member{Name: "Alex", Email: "alex@mail.com", Plan: ""},
			wantErr: ErrInvalidPlan,
		},
		{
			name:    "multiple errors returns first (email)",
			member:  Member{Name: "", Email: "", Plan: ""},
			wantErr: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.member.Validate()
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
