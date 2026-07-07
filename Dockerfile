# --- STAGE 1: BUILDER ---
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gym-server ./cmd/server/main.go

# --- STAGE 2: Final (Production) ---
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -h /home/gym -u 10001 gym

WORKDIR /home/gym

COPY --from=builder /app/gym-server .
COPY --from=builder /app/migrations ./migrations

USER gym

EXPOSE 50051

CMD ["./gym-server"]
