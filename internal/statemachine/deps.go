package statemachine

//go:generate mockgen -source=$GOFILE -destination=mocks/storage_mock.go

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
	"time"
)

// Storage интерфейс хранения данных стейтмашины
type Storage interface {
	RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error
	// CreateState создание нового стейта в базе
	CreateState(ctx context.Context, state *storage.State) error
	// GetStateByIdempotencyKey получение стейта по ключу idempotencyKey
	GetStateByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*storage.State, error)
	// GetStateByID получение стейта по id
	GetStateByID(ctx context.Context, stateID uuid.UUID) (*storage.State, error)
	// SaveStepExecuteInfo Сохранение информации о запуске выполнения шага
	SaveStepExecuteInfo(ctx context.Context, execute storage.StepExecuteInfo) error
	// UpdateState обновление стейта
	UpdateState(ctx context.Context, state storage.UpdateState) error
}

// UUIDGenerator интерфейс для генерации UUID (реальный или мок)
type UUIDGenerator interface {
	New() uuid.UUID
}

// Clock интерфейс для работы со временем (реальный или мок)
type Clock interface {
	Now() time.Time
}
