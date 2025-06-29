package statemachine

import (
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/statemachine/stepper"
	"time"
)

type Status = int

const (
	NewStatus        Status = iota
	InProgressStatus Status = iota
	CompletedStatus  Status = iota
	FailedStatus     Status = iota
)

// State состояние стейт машины
type State[DataT any, StatusT ~string, TypeT ~string] struct {
	// ID идентификаторе текущего стейта
	ID uuid.UUID
	// IdempotencyKey Ключ идемпотентности стейта
	IdempotencyKey uuid.UUID
	//  CreatedAt дата создания состояния
	CreatedAt time.Time
	// UpdatedAt дата обновления Status или Step
	UpdatedAt time.Time
	// Status статус
	Status Status
	// Step текущий шаг
	Step StatusT
	// Type тип состояния
	Type TypeT
	// Data данные выпуска
	Data DataT
}

// CreateState структура инициализации выпуска
type CreateState[DataT any] struct {
	// Data данные выпуска
	Data DataT
}

// UpdateState структура для обновление состояния стейт машины
type UpdateState[DataT any, StatusT ~string] struct {
	// UpdatedAt дата обновления Status или Step
	UpdatedAt time.Time
	// Status статус
	Status Status
	// Step текущий шаг
	Step StatusT
	// Data данные выпуска
	Data DataT
}

// StepExecuteInfo Информация о выполнении шагов стейт машины
type StepExecuteInfo[StatusT ~string] struct {
	StateID            uuid.UUID
	StartExecutedAt    time.Time
	CompleteExecutedAt time.Time
	Error              *string
	PreviewStep        StatusT
	NextStep           *StatusT
	NextStatus         Status
}

type CreateOptions interface {
	GetIdempotencyKey() uuid.UUID
}

type Step[DataT any, StatusT ~string, TypeT ~string] struct {
	OnStep stepper.StepFunc[DataT, StatusT, TypeT]
}

type StepRegistrationParams struct {
}

type StepRegistration[DataT any, StatusT ~string, TypeT ~string] struct {
	Steps map[StatusT]Step[DataT, StatusT, TypeT]
}
