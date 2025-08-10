package mkvmerge

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	PendingStatus  Status = iota
	RunningStatus  Status = iota
	CompleteStatus Status = iota
	ErrorStatus    Status = iota
)

type MergeLogs struct {
	CreatedAt time.Time
	Type      MessageType
	Content   string
}

type MergeResult struct {
	ID          uuid.UUID
	Params      MergeParams
	Status      Status
	Error       *string
	CreatedAt   time.Time
	CompletedAt *time.Time
	Progress    *float64
}

type CreateMergeResult struct {
	ID             uuid.UUID
	IdempotencyKey string
	Params         MergeParams
	Status         Status
	CreatedAt      time.Time
}

type UpdateMergeResult struct {
	Status    Status
	Error     *string
	Completed *time.Time
}

type MessageType int

const (
	InfoMessageType MessageType = iota
	ErrorMessageType
)

type OutputMessage struct {
	Type    MessageType
	Content string
}
