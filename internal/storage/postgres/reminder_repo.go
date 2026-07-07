package postgres

import (
	"context"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/service/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReminderRepo struct {
	pool *pgxpool.Pool
}

func NewReminderRepo(pool *pgxpool.Pool) *ReminderRepo {
	return &ReminderRepo{pool: pool}
}

// ListExpiringMemberships implements [service.ExpiringMembershipRepository].
// Joins members -> gyms -> users to find the gym owner's email for each
// expiring monthly membership. u.email is AES-encrypted at rest (see
// internal/utils.EncryptEmail) — decrypting it is the caller's job, not this
// repo's.
func (r *ReminderRepo) ListExpiringMemberships(ctx context.Context, within time.Duration) ([]*dto.MemberWithOwnerEmail, error) {
	threshold := time.Now().Add(within)
	query := `
		SELECT m.id, m.gym_id, m.name, m.phone, m.type, m.membership_start, m.membership_end,
		       u.email
		FROM members m
		JOIN gyms g ON g.id = m.gym_id
		JOIN users u ON u.id = g.owner_user_id
		WHERE m.type = 'monthly'
		  AND m.deleted_at IS NULL
		  AND m.membership_end IS NOT NULL
		  AND m.membership_end <= $1
		  AND m.membership_end >= NOW()`

	rows, err := r.pool.Query(ctx, query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*dto.MemberWithOwnerEmail
	for rows.Next() {
		var m domain.Member
		var memberType, ownerEmail string
		if err := rows.Scan(&m.ID, &m.GymID, &m.Name, &m.Phone, &memberType, &m.MembershipStart, &m.MembershipEnd, &ownerEmail); err != nil {
			return nil, err
		}
		m.Type = domainMemberType(memberType)
		results = append(results, &dto.MemberWithOwnerEmail{Member: &m, OwnerEmail: ownerEmail})
	}
	return results, rows.Err()
}
