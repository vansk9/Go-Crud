# Gunakan image golang untuk build
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Gunakan image minimal untuk menjalankan
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/app .
COPY .env .

EXPOSE 8080

CMD ["./app"]
