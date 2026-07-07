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
came from a real need — I approached a local gym with an early MVP,
which gave me a concrete problem to solve beyond a simple CRUD tutorial.
The current data model (monthly vs. occasional members, check-ins,
credit/"fiado" payment tracking, membership expiration reminders) was
shaped directly from conversations with real small-gym owners, not
guessed upfront.

**What this project demonstrates:**
- Hexagonal architecture with clear separation of domain, use cases, and adapters
- Multi-tenant data model: each gym is an isolated tenant, enforced at the JWT + query level
- Secure authentication: bcrypt (cost 14), AES-256-GCM email encryption, JWT HMAC-SHA256
- High-performance communication via gRPC with Protocol Buffers and server reflection
- Concurrent background processing using a worker pool pattern (goroutines + buffered channels)
- Full CI/CD pipeline: lint → test → build → Docker validation → migrations (GitHub Actions)
- Graceful shutdown handling SIGTERM/SIGINT with context cancellation

This project is a gym member management system designed under **Hexagonal Architecture** principles and high-performance communication using **gRPC**. It is built to be scalable, resilient, and highly concurrent.

---

## Architecture & Patterns

The project implements a clear separation of concerns:

- **Domain (`internal/domain`):** Pure business entities and rules — no dependency on protobuf, SQL, or any external package. `Member`, `Gym`, `CheckIn`, `User`, their validation, and domain-level errors live here.
- **Service (`internal/service`):** Use cases and orchestration — `MemberService`, `IdentityService`, `ReminderService`. Repository interfaces are defined here (inbound ports), never in `domain`.
- **DTOs (`internal/service/dto`):** Read-only types that cross aggregate boundaries (e.g. a member joined with its gym owner's email for reminders). Kept out of `domain` on purpose — they exist only to carry query results, not business rules.
- **Middleware (`internal/middleware`):** gRPC JWT interceptor — validates the token and injects the authenticated `gym_id`/`user_id` into `context.Context` for every handler.
- **Storage (`internal/storage/postgres`):** Persistence adapters using `pgx`, one file per aggregate (`member_repo.go`, `checkin_repo.go`, `identity_repo.go`, `reminder_repo.go`).
- **API/Proto (`api/proto`):** Contract definitions using Protocol Buffers, split by bounded context (`identity.proto`, `member.proto`, `checkin.proto`).
- **cmd/server:** The gRPC adapter and application entry point — `main.go` only wires dependencies and starts the server; handlers, proto↔domain mapping, and interceptors each live in their own file.
- **Concurrency Pattern:** Transport-agnostic **worker pool** (`internal/service/email_worker.go`) with buffered channels for background email processing, and **graceful shutdown** using `context.Context` + `sync.WaitGroup`.

### Multi-tenancy

- Each gym is an isolated tenant; a `gym_id` claim travels inside every JWT issued at login.
- A gRPC unary interceptor validates the token and injects `gym_id`/`user_id` into the request context on every call except `Login` and `RegisterUser`.
- Handlers and repository queries always scope by the `gym_id` from the authenticated context — it is **never** trusted from the request payload, preventing one gym from reading or modifying another gym's data.
- Registering an owner (`RegisterUser`) creates the `User` and their `Gym` atomically in a single database transaction — either both exist or neither does.

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
├── .github/workflows/ci.yml    # CI pipeline (lint, test, build, docker, migrations)
├── .golangci.yml                # Linter configuration
├── Makefile                     # Common development tasks
├── Dockerfile                   # Multi-stage build (non-root user, alpine:3.21)
├── docker-compose.yml           # DB + migrate + app orchestration
├── api/proto/                   # Protocol Buffer source definitions
│   ├── identity.proto           # Auth + User + atomic Gym registration contracts
│   ├── member.proto             # Member CRUD + check-in RPC contracts
│   └── checkin.proto            # Check-in / payment status messages
├── cmd/server/
│   ├── main.go                  # Bootstrap only: DB conn, wiring, graceful shutdown
│   ├── handlers.go              # gRPC handler methods (adapter layer)
│   ├── mappers.go               # proto <-> domain translation helpers
│   ├── interceptors.go          # logging + panic recovery interceptors
│   └── reminders.go             # daily membership-expiration reminder loop
├── internal/
│   ├── domain/                  # Pure business entities + validation (no external deps)
│   │   ├── member.go            # Member, MemberType, IsMembershipActive()
│   │   ├── user.go              # User
│   │   ├── gym.go                # Gym (tenant)
│   │   ├── checkin.go            # CheckIn, PaymentStatus
│   │   └── errors.go             # Sentinel domain errors
│   ├── middleware/
│   │   └── auth.go              # gRPC JWT interceptor; injects gym_id/user_id into context
│   ├── service/                 # Use cases, repository interfaces, worker pool
│   │   ├── member_service.go    # MemberService: CRUD + CheckIn + history
│   │   ├── identity_service.go  # IdentityService: atomic user+gym registration, login
│   │   ├── reminder_service.go  # ReminderService: daily expiring-membership job
│   │   ├── email_worker.go      # Transport-agnostic email worker pool
│   │   ├── dto/reporting.go     # MemberWithOwnerEmail (read-only cross-aggregate DTO)
│   │   └── repository_mock.go   # testify mocks for all repository interfaces
│   ├── storage/postgres/        # Output adapters (PostgreSQL via pgx)
│   │   ├── member_repo.go
│   │   ├── checkin_repo.go       # + checkin_repo_test.go (testcontainers)
│   │   ├── identity_repo.go     # includes RegisterUserWithGym (single tx)
│   │   ├── reminder_repo.go
│   │   └── test_setup.go        # shared integration test fixtures (setupTestDB, newTestGym, newTestMember)
│   ├── e2e/                      # full gRPC client-server integration tests — scaffolded, pending
│   └── utils/                    # Security (bcrypt, AES-GCM, JWT generate + validate)
├── migrations/                   # Database migrations (golang-migrate)
└── proto/                        # Generated .pb.go code
```

## gRPC Services

### GymService

| RPC | Request | Response | Description |
|-----|---------|----------|--------------|
| `CreateMember` | `MemberRequest` | `MemberResponse` | Register a new member (monthly or occasional). `gym_id` is taken from the authenticated context, never from the request |
| `ListMembers` | `ListMembersRequest` | `ListMembersResponse` | List all members for the authenticated gym |
| `UpdateMember` | `UpdateMemberRequest` | `MemberResponse` | Update a member's basic data and membership dates |
| `CheckIn` | `CheckInRequest` | `CheckInResponse` | Record a visit. Payment is only asked for occasional members — the server decides this from membership status, ignoring whatever the client sends for active monthly members |
| `GetMemberHistory` | `IdRequest` | `MemberHistoryResponse` | Full check-in history for a member — used as evidence to resolve payment disputes |

### UserService

| RPC | Request | Response | Description |
|-----|---------|----------|--------------|
| `RegisterUser` | `RegisterUserRequest` (includes `gym_name`) | `RegisterUserResponse` | Registers the owner and their gym atomically, in a single transaction |
| `GetUserProfile` | `UserRequest` | `User` | Get user profile |

### AuthService

| RPC | Request | Response | Description |
|-----|---------|----------|--------------|
| `Login` | `LoginRequest` | `LoginResponse` | Authenticate and get a JWT carrying `user_id`, `role`, and `gym_id` |
| `Logout` | `LogoutRequest` | `LogoutResponse` | Invalidate session |

---

## Security

- **Passwords:** Hashed with bcrypt (cost factor 14) via `golang.org/x/crypto`
- **Emails:** Encrypted with AES-256-GCM before storage
- **Authentication:** JWT tokens signed with HMAC-SHA256 (24h expiry), carrying `user_id`, `role`, and `gym_id`
- **Authorization:** gRPC interceptor validates every request except `Login`/`RegisterUser`; handlers trust only the `gym_id` from the verified token, never from request payloads
- **Keys:** `AES_MASTER_KEY` and `JWT_SECRET` are separate environment variables
- **Docker:** Runs as non-root user (uid 10001) on pinned `alpine:3.21`

## Concurrency

The project demonstrates a transport-agnostic **worker pool** pattern (`EmailWorkerPool`):

- 3 worker goroutines process email jobs (welcome emails, expiration reminders) from a buffered channel
- Workers respect context cancellation for ordered shutdown
- `sync.WaitGroup` tracks worker completion
- Non-blocking channel send (drops the job and logs a warning if the channel is full)
- A daily ticker (`cmd/server/reminders.go`) triggers `ReminderService.RunDaily`, which queries memberships expiring soon and enqueues an email to the gym owner

## Testing

| Layer | Tests | Framework |
|-------|-------|-----------|
| Domain | Member/User validation, `IsMembershipActive` | std testing |
| Service | `MemberService`, `IdentityService` (mocked repositories) | testify |
| Storage (Postgres) | Repository integration tests against a real Postgres instance | testcontainers-go |
| Utils | `HashPassword`, `EncryptEmail`, JWT generate + validate | testify |
| E2E | Full gRPC client-server tests via `bufconn` | scaffolded in `internal/e2e/`, pending |

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

- [x] Multi-tenant domain model: Gym, Member (monthly/occasional), CheckIn
- [x] JWT-based `gym_id` scoping via gRPC interceptor
- [x] Atomic owner + gym registration
- [x] Integration tests for the storage layer (testcontainers-go)
- [ ] E2E tests: full gRPC client-server integration tests (`internal/e2e`, currently scaffolded but empty)
- [ ] BFF Integration: connect the existing Next.js frontend to this API
- [ ] Deploy: Render (backend) + Vercel (frontend)
- [ ] Enhanced error handling: proper gRPC error codes (`AlreadyExists`, `NotFound`, etc.)
- [ ] API documentation auto-generated from proto files via `protoc-gen-doc`
