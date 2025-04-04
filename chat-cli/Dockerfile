FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/chat-cli

RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-cli ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /chat-cli .

CMD ["./chat-cli"]