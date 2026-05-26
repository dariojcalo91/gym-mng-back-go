# AGENTS.md - Gym Manager Backend (Go)

## Project Overview
Backend gRPC (`:50051`) for gym member management with hexagonal architecture.
Stack: Go 1.25, gRPC, PostgreSQL 15, pgx, testcontainers-go, testify.

## Repository
- Path: `/Users/dev/Development/gym-mng-back-go`
- Module: `github.com/dariojcalo91/gym-backend-go-ver`
- Remote: `origin` → `github.com/dariojcalo91/gym-mng-back-go`

## Session History (Day 1 — May 24, 2026)

### 18 Commits Realizados (master ahead of origin/master by 18)

```
ac6d671 docs: update readme with architecture, API reference, and changes summary
ddd682f feat(db): add updated_at trigger migration and CI migration validation
6dd0f99 ci: add GitHub Actions pipeline and golangci-lint configuration
e467975 test: add service and utils tests with TestMain env setup
5e37581 test(domain): add table-driven tests for Member and User validation
d18bf5d feat(server): add structured logging, recovery interceptor, health check, and server reflection
16aec88 chore(makefile): add Makefile with common development tasks
ac23bc2 docker: improve security with non-root user and pinned base image
8c86bdc feat(domain): add plan validation for members
04479ce security(config): remove hardcoded database fallback URL
97251cf docs(config): complete .env.example with all required variables
679a64e feat(api): implement GetMemberStatus and Logout handlers
ff8be0d security(keys): separate JWT secret from AES master key
77bae30 feat(server): add graceful shutdown with context-aware worker pool
60b0b0d refactor(service): rename package from gym to service
6dc0070 feat(domain): move domain errors from utils to internal/domain package
112120f test(worker): add worker pool concurrency tests
70c9b51 test(e2e): add gRPC end-to-end tests with bufconn in-memory server
d9c591f test(integration): add storage integration tests with testcontainers-go
```

### What Was Implemented

#### Architecture & Code Quality (6 commits)
- Domain errors moved from `internal/utils` → `internal/domain/errors.go`
- Package renamed: `package gym` → `package service`
- Graceful shutdown with SIGTERM/SIGINT + context-aware workers + WaitGroup
- JWT_SECRET separated from AES_MASTER_KEY
- Plan validation added to Member entity
- Hardcoded DB fallback URL removed
- Non-blocking channel send (drops email if channel full)

#### gRPC & Server (3 commits)
- `GetMemberStatus` and `Logout` handlers implemented
- Structured logging (`log/slog`) replacing std `log`
- Panic recovery interceptor + logging interceptor
- gRPC Health Checking Protocol (`grpc.health.v1`)
- Server Reflection enabled (`grpcurl`/`grpcui` ready)

#### Infrastructure (4 commits)
- Docker: non-root user (uid 10001), pinned `alpine:3.21`
- `docker-compose.yml` passes `AES_MASTER_KEY` and `JWT_SECRET`
- Makefile with build, test, lint, proto, docker targets
- `.env.example` with all required variables
- `.gitignore` updated for coverage and docs artifacts
- Migration `000003`: `updated_at` trigger for members table

#### Testing — 52 test cases (from 7 original)

| Layer | File | Cases | Type |
|-------|------|-------|------|
| Domain | `internal/domain/member_test.go` | 7 | Unit (table-driven) |
| Domain | `internal/domain/user_test.go` | 5 | Unit (table-driven) |
| Service | `internal/service/service_test.go` | 15 | Unit (testify mocks) |
| Worker Pool | `internal/service/worker_pool_test.go` | 4 | Unit (concurrency) |
| Security | `internal/utils/security_test.go` | 8 | Unit |
| JWT | `internal/utils/jwt_manager_test.go` | 3 | Unit |
| Storage | `internal/storage/postgres/repository_test.go` | 5 | Integration (testcontainers) |
| E2E | `cmd/main_test.go` | 5 | E2E (bufconn gRPC) |

#### CI/CD (2 commits)
- `.github/workflows/ci.yml` with 5 jobs: lint, test, build, docker, migrations
- `.golangci.yml` with 9 linters

### Known Issues to Fix

#### 1. CI Pipeline Red — Test Job
Test job needs env vars to skip integration/e2e tests (no Docker in some runners).
**Fix:** Add to `.github/workflows/ci.yml` test step:
```yaml
env:
  AES_MASTER_KEY: abcdefghijklmnopqrstuvwxyz123456
  JWT_SECRET: abcdefghijklmnopqrstuvwxyz123456
  SKIP_INTEGRATION: "true"
  SKIP_E2E: "true"
```

#### 2. CI Pipeline Red — Migrations Job
The action `gui-baiao/golang-migrate-action@v1` may not exist or fails.
**Fix:** Replace with direct `golang-migrate/migrate` CLI in a step:
```yaml
- name: Install migrate CLI
  run: |
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate /usr/local/bin/
- name: Run migrations
  run: migrate -path migrations -database "$DATABASE_URL" up
  env:
    DATABASE_URL: postgres://postgres:testpass@postgres:5432/testdb?sslmode=disable
```

#### 3. golangci-lint Not Installed Locally
`make lint` fails with `golangci-lint: No such file or directory`.
**Fix (macOS):** `brew install golangci-lint`
**Fix (other):** `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

#### 4. Module Path Mismatch (Low Priority)
`go.mod` says `github.com/dariojcalo91/gym-backend-go-ver` but repo is `gym-mng-back-go`.
Fix requires updating all imports across the codebase.

### Architecture Reference

```
cmd/main.go (gRPC handler + interceptors + health + reflection)
  │
  ├── internal/service/ (use cases + repository interfaces + worker pool)
  │     ├── gym.go              (GymService, Repository interface)
  │     ├── identity_service.go (IdentityService, IdentityRepository)
  │     ├── service_test.go
  │     └── worker_pool_test.go
  │
  ├── internal/domain/ (entities + validation + domain errors)
  │     ├── member.go + member_test.go
  │     ├── user.go + user_test.go
  │     └── errors.go
  │
  ├── internal/storage/postgres/ (output adapters)
  │     ├── repository.go + repository_test.go (testcontainers)
  │     └── user_repo.go
  │
  └── internal/utils/ (cross-cutting)
        ├── security.go + security_test.go
        └── jwt_manager.go + jwt_manager_test.go
```

### Commands Quick Reference

```bash
# Run all tests (skip integration/e2e)
SKIP_INTEGRATION=true SKIP_E2E=true go test -count=1 ./...

# Run with integration (requires Docker running)
AES_MASTER_KEY=... JWT_SECRET=... go test -count=1 ./...

# Run specific test
go test -v -count=1 -run "TestE2E" ./cmd/...

# Lint
golangci-lint run ./...

# Build
go build -o bin/gym-server ./cmd/main.go

# Docker
docker compose up --build -d

# API exploration (requires server running)
grpcui -plaintext localhost:50051
```

### User Preferences
- Language: Spanish (conversation), English (code/comments/docs)
- Conventional commits (`type(scope): description`)
- Atomic commits
- Must confirm before pushing to remote
- Must update README.md with relevant changes
- TDD approach preferred
