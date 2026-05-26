package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Skipping integration test: SKIP_INTEGRATION=true")
	}

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategyAndDeadline(
			60*time.Second,
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	runMigrations(t, pool)

	cleanup := func() {
		pool.Close()
		_ = pgContainer.Terminate(ctx)
	}

	return pool, cleanup
}

func runMigrations(t *testing.T, pool *pgxpool.Pool) {
	migrationDir := "../../../migrations"
	entries, err := os.ReadDir(migrationDir)
	require.NoError(t, err)

	ctx := context.Background()
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" && isUpMigration(entry.Name()) {
			content, err := os.ReadFile(filepath.Join(migrationDir, entry.Name()))
			require.NoError(t, err)

			_, err = pool.Exec(ctx, string(content))
			require.NoError(t, err, "failed to run migration: %s", entry.Name())
		}
	}
}

func isUpMigration(name string) bool {
	return len(name) > 8 && name[len(name)-7:] == ".up.sql"
}

func TestStorage_SaveAndGetMember(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	storage := NewStorage(pool)

	t.Run("save and retrieve member", func(t *testing.T) {
		member := &domain.Member{
			Name:  "Alex",
			Email: "alex@example.com",
			Plan:  "Gold",
		}

		err := storage.SaveMember(ctx, member)
		require.NoError(t, err)
		assert.NotEmpty(t, member.ID)

		retrieved, err := storage.GetMemberByID(ctx, member.ID)
		require.NoError(t, err)
		assert.Equal(t, member.Name, retrieved.Name)
		assert.Equal(t, member.Email, retrieved.Email)
		assert.Equal(t, member.Plan, retrieved.Plan)
		assert.Equal(t, "active", retrieved.Status)
	})

	t.Run("get non-existent member returns error", func(t *testing.T) {
		_, err := storage.GetMemberByID(ctx, "00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		member1 := &domain.Member{Name: "Bob", Email: "bob@example.com", Plan: "VIP"}
		err := storage.SaveMember(ctx, member1)
		require.NoError(t, err)

		member2 := &domain.Member{Name: "Bob2", Email: "bob@example.com", Plan: "Basic"}
		err = storage.SaveMember(ctx, member2)
		assert.Error(t, err)
	})

	t.Run("soft delete hides member from GetMemberByID", func(t *testing.T) {
		member := &domain.Member{Name: "Soft", Email: "soft@example.com", Plan: "Monthly"}
		err := storage.SaveMember(ctx, member)
		require.NoError(t, err)

		_, err = storage.pool.Exec(ctx, "UPDATE members SET deleted_at = NOW() WHERE id = $1", member.ID)
		require.NoError(t, err)

		_, err = storage.GetMemberByID(ctx, member.ID)
		assert.Error(t, err)
	})
}

func TestStorage_GetMemberByID_NotFound(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	storage := NewStorage(pool)

	_, err := storage.GetMemberByID(ctx, "00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}
