.PHONY: build run test test-cover lint proto doc docker-up docker-build clean

# Build
build:
	go build -ldflags="-s -w" -o bin/gym-server ./cmd/main.go

run:
	go run ./cmd/main.go

# Test
test:
	go test -v -count=1 ./...

test-cover:
	go test -v -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	golangci-lint run ./...

# Protobuf generation
proto:
	protoc --go_out=. --go-grpc_out=. api/proto/member.proto
	protoc --go_out=. --go-grpc_out=. api/proto/identity.proto

# Documentation generation (requires protoc-gen-doc)
proto-doc:
	mkdir -p docs
	protoc --doc_out=./docs --doc_opt=html,api.html api/proto/*.proto
	protoc --doc_out=./docs --doc_opt=markdown,api.md api/proto/*.proto

# Docker
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-build:
	docker build -t gym-backend .

# Migrations (requires golang-migrate CLI)
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Clean
clean:
	rm -rf bin/ coverage.out coverage.html docs/
