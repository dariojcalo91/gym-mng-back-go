# --- STAGE 1: BUILDER ---
# ensure you have the same go version in your 
# local environment and in the Dockerfile to avoid compatibility issues
FROM golang:1.24-alpine AS builder

# Install git and certificates (needed to download dependencies)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Compile binary:
# - CGO_ENABLED=0: Creates a static binary (does not depend on C libraries)
# - ldflags="-s -w": Reduces the size of the binary by removing debug symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gym-server ./cmd/main.go

# --- STAGE 2: Final (Production) ---
# We use 'scratch' or 'alpine' for the final image. Alpine is more friendly for debugging.
FROM alpine:latest

# Important for HTTPS connections and time zones
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Copy only the binary from the builder stage
COPY --from=builder /app/gym-server .

# Copy migrations (needed for the migrate service)
COPY --from=builder /app/migrations ./migrations

# Expose the gRPC port
EXPOSE 50051

# Ejecutamos el binario
CMD ["./gym-server"]
