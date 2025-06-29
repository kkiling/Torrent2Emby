package stepper

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/statemachine"
	"github.com/kkiling/torrent2emby/internal/storage"
	"github.com/samber/lo"
)

// Storage интерфейс хранения данных
type storageService[DataT any, StatusT ~string, TypeT ~string] interface {
	storage.Transactor
	// SaveStepExecuteInfo Сохранение информации о запуске выполнения шага
	SaveStepExecuteInfo(ctx context.Context, info statemachine.StepExecuteInfo[StatusT]) error
	// UpdateState обновление стейта
	UpdateState(ctx context.Context, state statemachine.UpdateState[DataT, StatusT]) error
}

type step[DataT any, StatusT ~string, TypeT ~string] struct {
	status StatusT
	onStep StepFunc[DataT, StatusT, TypeT]
}

// Stepper выполняет шаги стейт машины
type Stepper[DataT any, StatusT ~string, TypeT ~string] struct {
	storage storageService[DataT, StatusT, TypeT]
	steps   map[StatusT]step[DataT, StatusT, TypeT]
	clock   clock.Clock
}

// Add добавляет новый шаг в степпер
func (s *Stepper[DataT, StatusT, TypeT]) Add(status StatusT, onStep StepFunc[DataT, StatusT, TypeT]) {
	_, ok := s.steps[status]
	if ok {
		panic(fmt.Sprintf("steps already contains step: %s", status))
	}

	s.steps[status] = step[DataT, StatusT, TypeT]{
		status: status,
		onStep: onStep,
	}
}

// Compete выполняет стейт машину
func (s *Stepper[DataT, StatusT, TypeT]) Compete(
	ctx context.Context,
	inputState statemachine.State[DataT, StatusT, TypeT],
) (*statemachine.State[DataT, StatusT, TypeT], error) {

	if inputState.Status == statemachine.FailedStatus || inputState.Status == statemachine.CompletedStatus {
		return nil, fmt.Errorf("state already in terminal status")
	}
	currentState := inputState

	// Крутим стейт машину
	for ctx.Err() == nil {
		stepInfo, ok := s.steps[currentState.Step]
		if !ok {
			return nil, fmt.Errorf("unknown step %s", currentState.Step)
		}

		stepCtx := StepContext[DataT, StatusT, TypeT]{
			Data:  currentState.Data,
			State: currentState,
		}

		stepExecuteInfo := statemachine.StepExecuteInfo[StatusT]{
			StateID:         uuid.UUID{},
			StartExecutedAt: s.clock.Now(),
			PreviewStep:     currentState.Step,
		}

		stepResult := stepInfo.onStep(ctx, stepCtx)
		stepExecuteInfo.CompleteExecutedAt = s.clock.Now()

		newState := statemachine.State[DataT, StatusT, TypeT]{
			ID:        currentState.ID,
			CreatedAt: currentState.CreatedAt,
			UpdatedAt: currentState.UpdatedAt,
			Status:    currentState.Status,
			Step:      currentState.Step,
			Type:      currentState.Type,
			Data:      currentState.Data,
		}

		if stepResult.err != nil {
			stepExecuteInfo.Error = lo.ToPtr(stepResult.err.Error())
		}

		if stepResult.newData != nil {
			// Обновляем данные выпуска
			newState.Data = *stepResult.newData
		}

		isBreak := false
		switch stepResult.state {
		case emptyStepState:
			// Шаг не двигаем
			isBreak = true
		case nextStepState:
			if newState.Status == statemachine.NewStatus {
				newState.Status = statemachine.InProgressStatus
			}
			if newState.Step == *stepResult.nextStatus {
				return &newState, fmt.Errorf("error change to the same status")
			}
			newState.Step = *stepResult.nextStatus
			newState.UpdatedAt = s.clock.Now()
			stepExecuteInfo.NextStep = stepResult.nextStatus
		case failStepState:
			newState.Status = statemachine.FailedStatus
			newState.UpdatedAt = s.clock.Now()
			isBreak = true
		case completeStepState:
			newState.Status = statemachine.CompletedStatus
			newState.UpdatedAt = s.clock.Now()
			isBreak = true
		}
		stepExecuteInfo.NextStatus = newState.Status

		// TODO: Транзакция
		err := s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
			err := s.storage.SaveStepExecuteInfo(ctx, stepExecuteInfo)
			if err != nil {
				return fmt.Errorf("storage.SaveStepExecuteInfo: %w", err)
			}
			err = s.storage.UpdateState(ctx, statemachine.UpdateState[DataT, StatusT]{
				UpdatedAt: newState.UpdatedAt,
				Status:    newState.Status,
				Step:      newState.Step,
				Data:      newState.Data,
			})
			if err != nil {
				return fmt.Errorf("storage.UpdateState: %w", err)
			}
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("storage.RunTransaction: %w", err)
		}

		if isBreak {
			return &newState, nil
		}
		currentState = newState
	}

	return &currentState, nil
}
