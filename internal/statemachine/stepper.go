package statemachine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

// Stepper выполняет шаги стейт машины
type Stepper[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string] struct {
	storage Storage
	clock   Clock
	steps   map[StepT]Step[DataT, FailDataT, MetaDataT, StepT, TypeT]
}

func NewStepper[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string](
	storage Storage,
	clock Clock,
) *Stepper[DataT, FailDataT, MetaDataT, StepT, TypeT] {
	return &Stepper[DataT, FailDataT, MetaDataT, StepT, TypeT]{
		storage: storage,
		clock:   clock,
		steps:   make(map[StepT]Step[DataT, FailDataT, MetaDataT, StepT, TypeT]),
	}
}

// Add добавляет новый шаг в степпер
func (s *Stepper[DataT, FailDataT, MetaDataT, StepT, TypeT]) Add(status StepT, step Step[DataT, FailDataT, MetaDataT, StepT, TypeT]) {
	_, ok := s.steps[status]
	if ok {
		panic(fmt.Sprintf("steps already contains step: %s", status))
	}

	s.steps[status] = step
}

// Compete выполняет стейт машину
func (s *Stepper[DataT, FailDataT, MetaDataT, StepT, TypeT]) Compete(
	ctx context.Context,
	inputState State[DataT, FailDataT, MetaDataT, StepT, TypeT],
	options ...any,
) (st *State[DataT, FailDataT, MetaDataT, StepT, TypeT], executeErr error, err error) {
	if len(options) > 1 {
		return nil, nil, fmt.Errorf("too many options")
	}

	completeOptions := func() any {
		if len(options) == 1 {
			return options[0]
		}
		return nil
	}()

	currentState := inputState

	// Крутим стейт машину
	for ctx.Err() == nil {
		stepInfo, ok := s.steps[currentState.Step]
		if !ok {
			return nil, nil, fmt.Errorf("unknown step %s", currentState.Step)
		}

		execute := storage.StepExecuteInfo{
			StateID: currentState.ID,
			// Фиксация времени начала выполнения шага
			StartExecutedAt: s.clock.Now(),
			// Фиксация пред идущего шага
			PreviewStep: string(currentState.Step),
		}

		if stepInfo.OptionsType == nil && completeOptions != nil {
			return nil, nil, fmt.Errorf("optionsType is undefined and completeOptions is specified")
		}

		// Выполнение шага
		stepCtx := StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]{
			State:               currentState,
			completeOptionsType: stepInfo.OptionsType,
			completeOptions:     completeOptions,
		}

		stepResult := stepInfo.OnStep(ctx, stepCtx)
		newState := currentState

		// Фиксация времени выполнения шага
		execute.CompleteExecutedAt = s.clock.Now()
		if stepResult.newData != nil {
			// Обновляем данные выпуска
			newState.Data = *stepResult.newData
		}

		// Обработка
		isBreak := false
		switch stepResult.state {
		case emptyStepState:
			// Шаг не двигаем
			isBreak = true
		case errorStepState:
			// Сохранение ошибки выполнения шага если была
			execute.Error = lo.ToPtr(stepResult.err.Error())
			// Шаг не двигаем
			isBreak = true
		case nextStepState:
			if newState.Status == NewStatus {
				newState.Status = InProgressStatus
			}
			if newState.Step == *stepResult.nextStatus {
				return &newState, fmt.Errorf("error change to the same status"), nil
			}
			newState.Step = *stepResult.nextStatus
			newState.UpdatedAt = execute.CompleteExecutedAt
			execute.NextStep = lo.Ternary(stepResult.nextStatus != nil,
				lo.ToPtr(string(*stepResult.nextStatus)), nil)
		case failStepState:
			newState.Status = FailedStatus
			newState.UpdatedAt = execute.CompleteExecutedAt
			newState.Step = ""
			isBreak = true
		case completeStepState:
			newState.Status = CompletedStatus
			newState.UpdatedAt = execute.CompleteExecutedAt
			newState.Step = ""
			isBreak = true
		}

		// TODO: вынести в отделную функцию
		var data []byte
		data, err = json.Marshal(newState.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal data: %w", err)
		}
		if string(data) == "null" {
			data = []byte{}
		}

		var failData []byte
		failData, err = json.Marshal(newState.FailData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal failData: %w", err)
		}
		if string(failData) == "null" {
			failData = []byte{}
		}

		var metaData []byte
		metaData, err = json.Marshal(newState.MetaData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal metaData: %w", err)
		}
		if string(metaData) == "null" {
			metaData = []byte{}
		}
		err = s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
			terr := s.storage.SaveStepExecuteInfo(ctx, execute)
			if terr != nil {
				return fmt.Errorf("storage.SaveStepExecuteInfo: %w", terr)
			}

			terr = s.storage.UpdateState(ctx, newState.ID, storage.UpdateState{
				UpdatedAt: newState.UpdatedAt,
				Status:    newState.Status,
				Step:      string(newState.Step),
				Data:      data,
				FailData:  failData,
				MetaData:  metaData,
			})
			if terr != nil {
				return fmt.Errorf("storage.UpdateState: %w", terr)
			}
			return nil
		})

		if err != nil {
			return nil, nil, fmt.Errorf("storage.RunTransaction: %w", err)
		}

		if isBreak {
			// Возвращаем ошибку которую получили во время выполнения шага
			return &newState, stepResult.err, nil
		}

		currentState = newState
		// Сбрассываем опции, так как они нужны только для выполнения первого шага
		completeOptions = nil
	}

	return &currentState, nil, nil
}
