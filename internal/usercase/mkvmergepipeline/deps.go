package mkvmergepipeline

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
)

type MkvMerge interface {
	Merge(ctx context.Context, params mkvmerge.MergeParams, outputChan chan<- mkvmerge.OutputMessage) error
}

type Storage interface {
	Create(ctx context.Context, create *CreateMergeResult) error
	Update(ctx context.Context, id uuid.UUID, update *UpdateMergeResult) error
	DeleteLogs(ctx context.Context, mergeID uuid.UUID) error
	GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*MergeResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*MergeResult, error)
	GetOldestUncompleted(ctx context.Context) (*MergeResult, error)
	AddMergeLogs(ctx context.Context, id uuid.UUID, log MergeLogs) error
}
