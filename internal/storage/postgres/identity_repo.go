package postgres

import (
	"context"
	"fmt"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityRepo struct {
	pool *pgxpool.Pool
}

func NewIdentityRepo(pool *pgxpool.Pool) *IdentityRepo {
	return &IdentityRepo{pool: pool}
}

// GetUserByUsername implements [service.IdentityRepository].
func (r *IdentityRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	query := `SELECT id, gym_id, username, password, email, role FROM users WHERE username = $1`
	err := r.pool.QueryRow(ctx, query, username).Scan(&u.ID, &u.GymID, &u.Username, &u.Password, &u.Email, &u.Role)
	return u, err
}

// RegisterUserWithGym implements [service.IdentityRepository]. Creates the
// user and gym in one transaction — either both rows exist or neither does.
// There's a circular reference (gyms.owner_user_id -> users.id, users.gym_id
// -> gyms.id): resolved by inserting the user first without gym_id, then the
// gym, then backfilling users.gym_id, all inside the same tx.
func (r *IdentityRepo) RegisterUserWithGym(ctx context.Context, u *domain.User, gymName string) (*domain.Gym, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		u.Username, u.Password, u.Email, u.Role,
	).Scan(&u.ID)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	gym := &domain.Gym{Name: gymName, OwnerUserID: u.ID}
	err = tx.QueryRow(ctx,
		`INSERT INTO gyms (name, owner_user_id) VALUES ($1, $2) RETURNING id, created_at`,
		gym.Name, gym.OwnerUserID,
	).Scan(&gym.ID, &gym.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert gym: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET gym_id = $1 WHERE id = $2`, gym.ID, u.ID); err != nil {
		return nil, fmt.Errorf("link user to gym: %w", err)
	}
	u.GymID = gym.ID

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return gym, nil
}
