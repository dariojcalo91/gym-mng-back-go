package postgres

import (
	"context"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// GetUserByUsername implements [service.IdentityRepository].
func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	query := `SELECT id, username, password, email, role FROM users WHERE username = $1`
	err := r.pool.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Role)
	return u, err
}

// SaveUser implements [service.IdentityRepository].
func (r *UserRepo) SaveUser(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(ctx, query, u.Username, u.Password, u.Email, u.Role)
	if err != nil {
		return err
	}
	return nil
}
