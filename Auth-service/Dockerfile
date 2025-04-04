FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/Auth-service

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /auth-service .

CMD ["./auth-service"]