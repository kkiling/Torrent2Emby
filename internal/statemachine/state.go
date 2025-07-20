package statemachine

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Status = uint8

const (
	NewStatus        Status = iota
	InProgressStatus Status = iota
	CompletedStatus  Status = iota
	FailedStatus     Status = iota
)

// State состояние стейт машины
type State[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string] struct {
	// ID идентификаторе текущего стейта
	ID uuid.UUID
	// IdempotencyKey Ключ идемпотентности стейта
	IdempotencyKey string
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
	// FailData Данные фейла стейта
	FailData FailDataT
	// MetaDataT методаные стейта
	MetaData MetaDataT
}

// CreateState структура инициализации выпуска
type CreateState[DataT any, MetaDataT any, StepT ~string] struct {
	// FirstStep первый тип шага с которого начинать выполнение стейт машины
	FirstStep StepT
	// Data данные стейта
	Data DataT
	// MetaDataT методаные стейта
	MetaData MetaDataT
}

type CreateOptions interface {
	// GetIdempotencyKey Ключ идемпотентности
	GetIdempotencyKey() string
}

type Step[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string] struct {
	OptionsType reflect.Type
	OnStep      StepFunc[DataT, FailDataT, MetaDataT, StepT, TypeT]
}

type StepRegistrationParams struct {
}

type StepRegistration[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string] struct {
	Steps map[StepT]Step[DataT, FailDataT, MetaDataT, StepT, TypeT]
}
