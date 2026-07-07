package domain

import (
	"testing"
	"time"
)

func TestMember_Validate(t *testing.T) {
	tests := []struct {
		name    string
		member  Member
		wantErr error
	}{
		{
			name:    "valid member",
			member:  Member{Name: "Alex", Phone: "1234567890", Type: "MONTHLY", GymID: "Gym_1", MembershipStart: &time.Time{}, MembershipEnd: &time.Time{}},
			wantErr: nil,
		},
		{
			name:    "empty phone",
			member:  Member{Name: "Alex", Phone: "", Type: "OCCASIONAL"},
			wantErr: ErrInvalidPhone,
		},
		{
			name:    "wrong type",
			member:  Member{Name: "Alex", Phone: "1234567890", Type: "wrong type", GymID: "Gym_1", MembershipStart: &time.Time{}, MembershipEnd: &time.Time{}},
			wantErr: ErrInvalidMemberType,
		},
		{
			name:    "whitespace name",
			member:  Member{Name: "  ", Phone: "alex@mail.com", Type: "MONTHLY"},
			wantErr: ErrInvalidName,
		},
		{
			name:    "empty type",
			member:  Member{Name: "Alex", Phone: "alex@mail.com", Type: "", GymID: "Gym_1", MembershipStart: &time.Time{}, MembershipEnd: &time.Time{}},
			wantErr: ErrInvalidMemberType,
		},
		{
			name:    "No gym",
			member:  Member{Name: "Alex", Phone: "alex@mail.com", Type: ""},
			wantErr: ErrInvalidGym,
		},
		{
			name:    "No membershipt start time",
			member:  Member{Name: "Alex", Phone: "alex@mail.com", Type: "MONTHLY", GymID: "Gym_1"},
			wantErr: ErrInvalidMembershipDates,
		},
		{
			name:    "No membershipt end time",
			member:  Member{Name: "Alex", Phone: "alex@mail.com", Type: "MONTHLY", GymID: "Gym_1", MembershipStart: &time.Time{}},
			wantErr: ErrInvalidMembershipDates,
		},
		{
			name:    "multiple errors returns first (name)",
			member:  Member{Name: "", Phone: "", Type: ""},
			wantErr: ErrInvalidName,
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
