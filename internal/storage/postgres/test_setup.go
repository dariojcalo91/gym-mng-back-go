package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Skipping checkin repo integration test: SKIP_INTEGRATION=true")
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

func isUpMigration(name string) bool {
	return len(name) > 8 && name[len(name)-7:] == ".up.sql"
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
