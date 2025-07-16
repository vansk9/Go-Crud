# Stage 1 - Build
FROM golang:1.23-alpine AS builder


WORKDIR /app

# Copy go mod & download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build dari file yang ada di folder cmd
RUN go build -o main ./cmd/main.go

# Stage 2 - Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main", "serve"]
