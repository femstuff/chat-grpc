FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/Saga-orchestrator

RUN CGO_ENABLED=0 GOOS=linux go build -o /saga-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /saga-service .

CMD ["./saga-service"]