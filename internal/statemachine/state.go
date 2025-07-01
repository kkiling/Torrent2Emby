package statemachine

import (
	"github.com/google/uuid"
	"reflect"
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
type State[DataT any, StepT ~string, TypeT ~string] struct {
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
	Step StepT
	// Type тип состояния
	Type TypeT
	// Data данные выпуска
	Data DataT
}

// CreateState структура инициализации выпуска
type CreateState[DataT any, StepT ~string] struct {
	// FirstStep первый тип шага с которого начинать выполнение стейт машины
	FirstStep StepT
	// Data данные выпуска
	Data DataT
}

type CreateOptions interface {
	GetIdempotencyKey() uuid.UUID
}

type Step[DataT any, StepT ~string, TypeT ~string] struct {
	OptionsType reflect.Type
	OnStep      StepFunc[DataT, StepT, TypeT]
}

type StepRegistrationParams struct {
}

type StepRegistration[DataT any, StepT ~string, TypeT ~string] struct {
	Steps map[StepT]Step[DataT, StepT, TypeT]
}
