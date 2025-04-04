FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/Chat-service

RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /chat-service .

CMD ["./chat-service"]