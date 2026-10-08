# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kv-server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kv-server .

RUN mkdir -p /app/data

EXPOSE 50051

CMD ["./kv-server", "-id", "node1", "-port", "50051", "-data", "/app/data"]
