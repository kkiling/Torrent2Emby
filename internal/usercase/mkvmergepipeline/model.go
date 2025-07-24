package mkvmergepipeline

import (
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	"time"
)

type Status int

const (
	Pending  Status = iota
	Running  Status = iota
	Complete Status = iota
	Error    Status = iota
)

type MergeLogs struct {
	CreatedAt time.Time
	Type      mkvmerge.MessageType
	Content   string
}

type MergeResult struct {
	ID          uuid.UUID
	Params      mkvmerge.MergeParams
	Status      Status
	Error       *string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type CreateMergeResult struct {
	ID             uuid.UUID
	IdempotencyKey string
	Params         mkvmerge.MergeParams
	Status         Status
	CreatedAt      time.Time
}

type UpdateMergeResult struct {
	Status    Status
	Error     *string
	Completed *time.Time
}
