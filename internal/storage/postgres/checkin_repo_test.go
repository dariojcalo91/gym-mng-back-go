package postgres

import (
	"context"
	"testing"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_SaveCheckIn(t *testing.T) {
	// Setup test database
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	checkinStorage := NewCheckInRepo(pool)
	identityStorage := NewIdentityRepo(pool)
	memberStorage := NewMemberRepo(pool)

	t.Run("Save checkin and retrieve list of members", func(t *testing.T) {
		user := &domain.User{
			Username: "jessy_gym_user",
		}

		gym, err := identityStorage.RegisterUserWithGym(ctx, user, "Test gym")
		require.NoError(t, err)

		member := &domain.Member{
			Name:  "Martin",
			Type:  domain.MemberTypeMonthly,
			GymID: gym.ID,
			Phone: "1234567890",
		}

		err = memberStorage.SaveMember(ctx, member)
		require.NoError(t, err)

		checkin := &domain.CheckIn{
			GymID:         user.GymID,
			MemberID:      member.ID,
			PaymentStatus: domain.PaymentStatusPaid,
		}

		err = checkinStorage.SaveCheckIn(ctx, checkin)
		require.NoError(t, err)
		assert.NotEmpty(t, checkin.ID)

		retrieved, err := checkinStorage.ListCheckInsByMember(ctx, checkin.GymID, checkin.MemberID)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, checkin.ID, retrieved[0].ID)
		assert.Equal(t, checkin.GymID, retrieved[0].GymID)
		assert.Equal(t, checkin.MemberID, retrieved[0].MemberID)
	})
}
