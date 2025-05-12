package entity

import "google.golang.org/protobuf/types/known/timestamppb"

type SagaStatus string

const (
	StatusPending   SagaStatus = "pending"
	StatusCompleted SagaStatus = "completed"
	SttusFaied      SagaStatus = "failed"
)

type Saga struct {
	ID         string
	MessageID  int64
	Status     SagaStatus
	Retries    int
	MaxRetries int
	CreatedAt  timestamppb.Timestamp
}
