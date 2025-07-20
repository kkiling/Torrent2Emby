package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotFound объект не найден в базе
	ErrNotFound = errors.New("entity not found")
	// ErrAlreadyExists запись уже существует
	ErrAlreadyExists = errors.New("entity already exists")
)

// State состояние стейт машины
type State struct {
	// ID идентификаторе текущего стейта
	ID uuid.UUID
	// IdempotencyKey Ключ идемпотентности стейта
	IdempotencyKey string
	//  CreatedAt дата создания состояния
	CreatedAt time.Time
	// UpdatedAt дата обновления Status или Step
	UpdatedAt time.Time
	// Status статус
	Status uint8
	// Step текущий шаг
	Step string
	// Type тип состояния
	Type string
	// Data данные выпуска
	Data []byte
	// FailData Данные фейла стейта
	FailData []byte
	// MetaDataT методаные стейта
	MetaData []byte
}

// UpdateState структура для обновление состояния стейт машины
type UpdateState struct {
	// UpdatedAt дата обновления Status или Step
	UpdatedAt time.Time
	// Status статус
	Status uint8
	// Step текущий шаг
	Step string
	// Data данные выпуска
	Data []byte
	// FailData Данные фейла стейта
	FailData []byte
	// MetaDataT методаные стейта
	MetaData []byte
}

// StepExecuteInfo Информация о выполнении шагов стейт машины
type StepExecuteInfo struct {
	StateID            uuid.UUID
	StartExecutedAt    time.Time
	CompleteExecutedAt time.Time
	Error              *string
	PreviewStep        string
	NextStep           *string
}
