# 🏋️ GYM Backend Manager

> A production-grade backend built with Go to manage gym members — 
> designed as a personal portfolio project to demonstrate real-world 
> engineering practices: hexagonal architecture, gRPC, TDD, security, 
> concurrency, and CI/CD.

[![CI](https://github.com/dariojcalo91/gym-mng-back-go/actions/workflows/ci.yml/badge.svg)](https://github.com/dariojcalo91/gym-mng-back-go/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![gRPC](https://img.shields.io/badge/gRPC-Protocol%20Buffers-blueviolet)
![License](https://img.shields.io/badge/license-MIT-green)

## Why this project?

Built during a deliberate practice period to keep skills sharp and 
explore production-grade patterns in Go. The domain (gym management) 
came from a real need — I approached a local gym with this MVP, 
which gave me a concrete problem to solve beyond a simple CRUD tutorial.

**What this project demonstrates:**
- Hexagonal architecture with clear separation of domain, use cases, and adapters
- Secure authentication: bcrypt (cost 14), AES-256-GCM email encryption, JWT HMAC-SHA256
- High-performance communication via gRPC with Protocol Buffers and server reflection
- Concurrent background processing using a worker pool pattern (goroutines + buffered channels)
- Full CI/CD pipeline: lint → test → build → Docker validation → migrations (GitHub Actions)
- Graceful shutdown handling SIGTERM/SIGINT with context cancellation

This project is a gym member management system designed under **Hexagonal Architecture** principles and high-performance communication using **gRPC**. It is built to be scalable, resilient, and highly concurrent.

---

## Architecture & Patterns

The project implements a clear separation of concerns:

- **Domain:** Pure business logic and entities (Validations, errors).
- **Service (internal/service):** Use cases and orchestration of asynchronous processes.
- **Storage (internal/storage/postgres):** Persistence adapter using `pgx`.
- **API/Proto:** Contract definitions using Protocol Buffers.
- **Concurrency Pattern:** Implementation of **Worker Pools** with channels for background task processing (e.g., sending emails) and **Graceful Shutdown** using `context.Context` and `sync.WaitGroup`.

---

## Quick Start (Docker)

Copy the environment file and customize your credentials:

```bash
cp .env.example .env
```

The project is fully orchestrated. To start the database, apply migrations, and launch the server:

```bash
docker compose up --build
```

---

## Development Commands

A `Makefile` is provided for common tasks:

```bash
make build       # Build binary to bin/gym-server
make run         # Run the server directly
make test        # Run all tests
make test-cover  # Run tests with coverage report
make lint        # Run golangci-lint
make proto       # Regenerate protobuf code
make proto-doc   # Generate API documentation from .proto files
make docker-up   # Start all services with docker compose
make clean       # Remove build artifacts
```

Run tests:

```bash
go test -count=1 ./...
```

---

## API Documentation

### Server Reflection

The gRPC server has **reflection** enabled, allowing runtime API discovery with tools like `grpcurl` and `grpcui`:

```bash
# Install grpcui
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest

# Launch interactive web UI
grpcui -plaintext localhost:50051
```

### Health Check

The server implements the standard **gRPC Health Checking Protocol** (`grpc.health.v1.Health`). During graceful shutdown, the health status is set to `NOT_SERVING`.

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `AES_MASTER_KEY` | Yes | 32-byte key for AES-256-GCM email encryption |
| `JWT_SECRET` | Yes | Secret key for HMAC-SHA256 JWT signing |

See `.env.example` for the full list.

---

## Project Structure

```
.
├── .github/workflows/ci.yml   # CI pipeline (lint, test, build, docker, migrations)
├── .golangci.yml               # Linter configuration
├── Makefile                    # Common development tasks
├── Dockerfile                  # Multi-stage build (non-root user, alpine:3.21)
├── docker-compose.yml          # DB + migrate + app orchestration
├── api/proto/                  # Protocol Buffer source definitions
├── cmd/
│   ├── main.go                 # Server entry point with graceful shutdown
│   └── client_test/            # gRPC client test script (gitignored)
├── internal/
│   ├── domain/                 # Entities, validation, domain errors
│   ├── service/                # Use cases, repository interfaces, worker pool
│   ├── storage/postgres/       # Output adapters (PostgreSQL)
│   └── utils/                  # Security (bcrypt, AES-GCM, JWT)
├── migrations/                 # Database migrations (golang-migrate)
└── proto/                      # Generated .pb.go code
```

## gRPC Services

### GymService

| RPC | Request | Response | Description |
|-----|---------|----------|-------------|
| `CreateMember` | `MemberRequest` | `MemberResponse` | Register a new member |
| `GetMemberStatus` | `IdRequest` | `MemberResponse` | Get member status by ID |

### UserService

| RPC | Request | Response | Description |
|-----|---------|----------|-------------|
| `RegisterUser` | `RegisterUserRequest` | `RegisterUserResponse` | Register a new user |
| `GetUserProfile` | `UserRequest` | `User` | Get user profile |

### AuthService

| RPC | Request | Response | Description |
|-----|---------|----------|-------------|
| `Login` | `LoginRequest` | `LoginResponse` | Authenticate and get JWT token |
| `Logout` | `LogoutRequest` | `LogoutResponse` | Invalidate session |

---

## Security

- **Passwords:** Hashed with bcrypt (cost factor 14) via `golang.org/x/crypto`
- **Emails:** Encrypted with AES-256-GCM before storage
- **Authentication:** JWT tokens signed with HMAC-SHA256 (24h expiry)
- **Keys:** `AES_MASTER_KEY` and `JWT_SECRET` are separate environment variables
- **Docker:** Runs as non-root user (uid 10001) on pinned `alpine:3.21`

## Concurrency

The project demonstrates a **worker pool** pattern:

- 3 worker goroutines process welcome emails from a buffered channel (capacity 5)
- Workers respect context cancellation for ordered shutdown
- `sync.WaitGroup` tracks worker completion
- Non-blocking channel send (drops email if channel full)

## Testing

| Layer | Tests | Framework |
|-------|-------|-----------|
| Domain | Member/User validation (12 cases) | std testing |
| Service | RegisterMember, RegisterUser, LoginUser, GetMemberStatus, GetUserByUsername (15 cases) | testify |
| Utils | HashPassword, EncryptEmail, JWT generation (11 cases) | testify |

Test commands:

```bash
make test          # Run all tests
make test-cover    # With coverage report
```

## DevOps & Resilience

- **Healthchecks:** Postgres health check before migration and app start
- **Graceful Shutdown:** Captures SIGTERM/SIGINT, cancels workers via context, performs `grpc.GracefulStop()`, waits for all workers to finish
- **Recovery Interceptor:** Panic recovery middleware for all gRPC handlers
- **Logging Interceptor:** Structured request/error logging for all gRPC calls
- **CI Pipeline:** GitHub Actions runs lint, test, build, Docker validation, and migration validation on every push

---

## Roadmap

- [ ] BFF Integration: Create a Backend-for-Frontend using Next.js
- [ ] Web Dashboard: Build a React-based dashboard to manage members visually
- [ ] Integration Tests: Add testcontainers-go for storage layer testing
- [ ] E2E Tests: Full gRPC client-server integration tests
- [ ] Enhanced Error Handling: gRPC error codes (AlreadyExists, etc.)
- [ ] API Documentation: Auto-generated from proto files via `protoc-gen-doc`
