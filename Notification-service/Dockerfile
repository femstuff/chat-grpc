FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/Notification-service

RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /notification-service .

CMD ["./notification-service"]