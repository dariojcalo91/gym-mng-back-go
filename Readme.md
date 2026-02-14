# 🏋️ GYM Backend Manager (Hexagonal + gRPC + Go)

This project is a gym member management system designed under **Hexagonal Architecture** principles and high-performance communication using **gRPC**. It is built to be scalable, resilient, and highly concurrent.

---

## 🏗️ Architecture & Patterns
The project implements a clear separation of concerns:
* **Domain:** Pure business logic and entities (Validations, Enums).
* **Internal/Gym (Service):** Use cases and orchestration of asynchronous processes.
* **Internal/Storage (Postgres):** Persistence adapter using `pgx`.
* **API/Proto:** Contract definitions using Protocol Buffers.
* **Concurrency Pattern:** Implementation of **Worker Pools** with channels for background task processing (e.g., sending emails) and **Graceful Shutdown** using `sync.WaitGroup`.

---

## 🚀 Quick Start (Docker)

**Note:** I've included a ```.env.example``` file. To run the project, copy it to ```.env``` and customize your credentials.

The project is fully orchestrated. To start the database, apply migrations, and launch the server:

```bash
docker-compose up --build
```

### 🛠️ Development Notes (Cheat Sheet)
1. gRPC & Protobuf
We use gRPC for efficient communication. The contract is defined in .proto files.

Requirements:

Install compiler: ```brew install protobuf```

Go Plugins:

```Bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Critical Step:** Ensure Golang BIN folder is in your PATH:
```Bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Generate Code:

```Bash
protoc --go_out=. --go-grpc_out=. api/proto/member.proto
```
Note: Never edit .pb.go files manually.


2. Database Migrations
We handle DB schema versioning to ensure consistency.

Tool: golang-migrate (```brew install golang-migrate```)

Create new migration:
```bash
migrate create -ext sql -dir migrations -seq migration_name
```

Structure:

up: Apply changes.

down: Rollback changes.

3. Goroutines & Concurrency
Harnessing Go's power for I/O Bound tasks.

Channels: "Don't communicate by sharing memory, share memory by communicating."

WaitGroups: Coordinating to wait for background processes before shutdown.

Context: Timeout control and request cancellation.

Implementation:
We use a Worker Pool for sending welcome emails. This prevents the user from waiting for the SMTP server to respond before receiving their registration confirmation.

---

## 🧪 Testing (TDD)
To run unit tests and business logic validation:

```Bash
go test -v ./internal/gym/...
```

## 🛡️ DevOps & Resilience
Healthchecks: The application service waits for Postgres to be HEALTHY before attempting to connect.

Graceful Shutdown: Upon receiving a termination signal (SIGTERM), the system drains channels and waits for workers to finish pending tasks before exiting, preventing data loss.

## 📈 Next Steps (Roadmap)
[ ] BFF Integration: Create a Backend-for-Frontend using Next.js.

[ ] Web Dashboard: Build a React-based dashboard to manage members visually.

[ ] Enhanced Logging: Implement structured logging (e.g., slog or zap).

[ ] Error Handling: Implement detailed gRPC error codes (e.g., AlreadyExists for duplicate emails).
