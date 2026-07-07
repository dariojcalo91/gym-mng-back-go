package postgres

import (
	"context"
	"strings"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckInRepo struct {
	pool *pgxpool.Pool
}

func NewCheckInRepo(pool *pgxpool.Pool) *CheckInRepo {
	return &CheckInRepo{pool: pool}
}

func dbPaymentStatus(s domain.PaymentStatus) string {
	return strings.ToLower(string(s))
}

func domainPaymentStatus(s string) domain.PaymentStatus {
	return domain.PaymentStatus(strings.ToUpper(s))
}

// SaveCheckIn implements [service.CheckInRepository].
func (r *CheckInRepo) SaveCheckIn(ctx context.Context, c *domain.CheckIn) error {
	query := `
		INSERT INTO check_ins (gym_id, member_id, checked_in_at, payment_status, payment_note)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.pool.QueryRow(ctx, query,
		c.GymID, c.MemberID, c.CheckedInAt, dbPaymentStatus(c.PaymentStatus), c.PaymentNote,
	).Scan(&c.ID)
}

// ListCheckInsByMember implements [service.CheckInRepository].
func (r *CheckInRepo) ListCheckInsByMember(ctx context.Context, gymID, memberID string) ([]*domain.CheckIn, error) {
	query := `
		SELECT id, gym_id, member_id, checked_in_at, payment_status, payment_note
		FROM check_ins
		WHERE member_id = $1 AND gym_id = $2
		ORDER BY checked_in_at DESC`

	rows, err := r.pool.Query(ctx, query, memberID, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkIns []*domain.CheckIn
	for rows.Next() {
		var c domain.CheckIn
		var paymentStatus string
		if err := rows.Scan(&c.ID, &c.GymID, &c.MemberID, &c.CheckedInAt, &paymentStatus, &c.PaymentNote); err != nil {
			return nil, err
		}
		c.PaymentStatus = domainPaymentStatus(paymentStatus)
		checkIns = append(checkIns, &c)
	}
	return checkIns, rows.Err()
}
