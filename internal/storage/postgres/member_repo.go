package postgres

import (
	"context"
	"strings"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberRepo struct {
	pool *pgxpool.Pool
}

func NewMemberRepo(pool *pgxpool.Pool) *MemberRepo {
	return &MemberRepo{pool: pool}
}

// dbMemberType/domainMemberType translate between the Postgres enum
// ('monthly', lowercase) and the domain's Go constants (MONTHLY, uppercase).
// The domain owns its own representation; this boundary adapts it to storage.
func dbMemberType(t domain.MemberType) string {
	return strings.ToLower(string(t))
}

func domainMemberType(s string) domain.MemberType {
	return domain.MemberType(strings.ToUpper(s))
}

// SaveMember implements [service.MemberRepository].
func (r *MemberRepo) SaveMember(ctx context.Context, m *domain.Member) error {
	query := `
		INSERT INTO members (gym_id, name, phone, type, membership_start, membership_end)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.pool.QueryRow(ctx, query,
		m.GymID, m.Name, m.Phone, dbMemberType(m.Type), m.MembershipStart, m.MembershipEnd,
	).Scan(&m.ID)
}

// GetMemberByID implements [service.MemberRepository]. Scoped by gym_id so one
// gym can never read another gym's member, even guessing a valid UUID.
func (r *MemberRepo) GetMemberByID(ctx context.Context, gymID, id string) (*domain.Member, error) {
	var m domain.Member
	var memberType string
	query := `
		SELECT id, gym_id, name, phone, type, membership_start, membership_end
		FROM members
		WHERE id = $1 AND gym_id = $2 AND deleted_at IS NULL`

	err := r.pool.QueryRow(ctx, query, id, gymID).Scan(
		&m.ID, &m.GymID, &m.Name, &m.Phone, &memberType, &m.MembershipStart, &m.MembershipEnd,
	)
	if err != nil {
		return nil, err
	}
	m.Type = domainMemberType(memberType)
	return &m, nil
}

// ListMembers implements [service.MemberRepository].
func (r *MemberRepo) ListMembers(ctx context.Context, gymID string) ([]*domain.Member, error) {
	query := `
		SELECT id, gym_id, name, phone, type, membership_start, membership_end
		FROM members
		WHERE gym_id = $1 AND deleted_at IS NULL
		ORDER BY name`

	rows, err := r.pool.Query(ctx, query, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.Member
	for rows.Next() {
		var m domain.Member
		var memberType string
		if err := rows.Scan(&m.ID, &m.GymID, &m.Name, &m.Phone, &memberType, &m.MembershipStart, &m.MembershipEnd); err != nil {
			return nil, err
		}
		m.Type = domainMemberType(memberType)
		members = append(members, &m)
	}
	return members, rows.Err()
}

// UpdateMember implements [service.MemberRepository].
func (r *MemberRepo) UpdateMember(ctx context.Context, m *domain.Member) error {
	query := `
		UPDATE members
		SET name = $1, phone = $2, membership_start = $3, membership_end = $4
		WHERE id = $5 AND gym_id = $6 AND deleted_at IS NULL`

	tag, err := r.pool.Exec(ctx, query, m.Name, m.Phone, m.MembershipStart, m.MembershipEnd, m.ID, m.GymID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrMemberNotFound
	}
	return nil
}
