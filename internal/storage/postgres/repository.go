package postgres

import (
	"context"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage is the adapter for Postgres
type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

// SaveMember inserts a member into the real database
func (s *Storage) SaveMember(ctx context.Context, m *domain.Member) error {
	query := `INSERT INTO members (name, email, plan, status) VALUES ($1, $2, $3, $4)`

	_, err := s.pool.Exec(ctx, query, m.Name, m.Email, m.Plan, m.Status)
	if err != nil {
		return err
	}
	return nil
}
