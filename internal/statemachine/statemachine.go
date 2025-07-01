package statemachine

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

type Config struct {
}

type StateMachine[DataT any, StepT ~string, TypeT ~string, CreateOptionsT CreateOptions] struct {
	cfg           Config
	runner        Runner[DataT, StepT, TypeT, CreateOptionsT]
	storage       Storage
	clock         Clock
	uuidGenerator UUIDGenerator
}

func NewService[DataT any, StepT ~string, TypeT ~string, CreateOptionsT CreateOptions](
	cfg Config,
	storage Storage,
	runner Runner[DataT, StepT, TypeT, CreateOptionsT],
) *StateMachine[DataT, StepT, TypeT, CreateOptionsT] {

	sm := StateMachine[DataT, StepT, TypeT, CreateOptionsT]{
		cfg:           cfg,
		runner:        runner,
		storage:       storage,
		clock:         &realClock{},
		uuidGenerator: &uuidGenerator{},
	}

	return &sm
}

func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) getStateByIdempotencyKey(
	ctx context.Context,
	idempotencyKey uuid.UUID,
) (*State[DataT, StepT, TypeT], error) {
	findState, err := i.storage.GetStateByIdempotencyKey(ctx, idempotencyKey)
	switch {
	case err == nil: // Выпуск найден
		return mapStorageToState[DataT, StepT, TypeT](findState)
	case errors.Is(err, storage.ErrNotFound): // Выпуск не найден
		return nil, nil
	default:
		return nil, fmt.Errorf("storage.GetStateByIdempotencyKey: %w", err)
	}
}

func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) getStateByID(
	ctx context.Context,
	stateID uuid.UUID,
) (*State[DataT, StepT, TypeT], error) {
	findState, err := i.storage.GetStateByID(ctx, stateID)
	if err != nil {
		// TODO: вернуть бизнес ошибку notfound
		return nil, fmt.Errorf("i.storage.GetState: %w", err)
	}
	return mapStorageToState[DataT, StepT, TypeT](findState)
}

// Create создание стейт машины
func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) Create(
	ctx context.Context,
	options CreateOptionsT,
) (*State[DataT, StepT, TypeT], error) {
	// Проверяем выпуск на наличие ключа идемпотентности
	if findState, err := i.getStateByIdempotencyKey(ctx, options.GetIdempotencyKey()); err != nil {
		return nil, fmt.Errorf("getStateByIdempotencyKey: %w", err)
	} else if findState != nil {
		return findState, nil
	}

	now := i.clock.Now()
	// Вызываем раннер, который выполняет бизнес логику и возвращает issueData
	create, err := i.runner.Create(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("runner.Create: %w", err)
	}

	newIssue := State[DataT, StepT, TypeT]{
		ID:             i.uuidGenerator.New(),
		IdempotencyKey: options.GetIdempotencyKey(),
		CreatedAt:      now,
		UpdatedAt:      now,
		Status:         NewStatus,
		Step:           create.FirstStep,
		Type:           i.runner.Type(),
		Data:           create.Data,
	}

	newStorageState, err := mapStateToStorage[DataT, StepT, TypeT](&newIssue)
	if err != nil {
		return nil, fmt.Errorf("mapStateT2: %w", err)
	}

	saveErr := i.storage.CreateState(ctx, newStorageState)
	if saveErr != nil {
		return nil, fmt.Errorf("storage.CreateState: %w", saveErr)
	}

	return &newIssue, nil
}

func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) initStepper() *Stepper[DataT, StepT, TypeT] {
	stepper := NewStepper[DataT, StepT, TypeT](i.storage, i.clock)
	stepsRegistration := i.runner.StepRegistration(StepRegistrationParams{})
	for s, step := range stepsRegistration.Steps {
		stepper.Add(s, step)
	}
	return stepper
}

// Complete выполнение выпуска
func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) Complete(
	ctx context.Context,
	stateID uuid.UUID,
	options ...any,
) (st *State[DataT, StepT, TypeT], executeErr error, err error) {
	// Проверяем выпуск на наличие ключа идемпотентности
	findState, err := i.getStateByID(ctx, stateID)
	if err != nil {
		return nil, nil, fmt.Errorf("getStateByID: %w", err)
	}

	stepper := i.initStepper()
	res, eErr, err := stepper.Compete(ctx, *findState, options...)
	if err != nil {
		return nil, nil, fmt.Errorf("stepper.Compete: %w", err)
	}

	return res, eErr, nil
}

// SetClock устанавливает кастомную реализацию часов
func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) SetClock(clock Clock) {
	i.clock = clock
}

// SetUUIDGenerator устанавливает кастомную реализацию uuid генератора
func (i *StateMachine[DataT, StepT, TypeT, CreateOptionsT]) SetUUIDGenerator(generator UUIDGenerator) {
	i.uuidGenerator = generator
}
