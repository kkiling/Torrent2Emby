package teststate

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/statemachine"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func mustBeData(data Data) []byte {
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return d
}

func TestTaskRunner_StepRegistration(t *testing.T) {
	t.Parallel()
	var (
		createdAt  = time.Now()
		stateID    = uuid.New()
		createOpts = CreateOptions{
			IdempotencyKey: uuid.New(),
			Title:          "Custom title",
			Amount:         42,
		}
		deps = setupTestDeps(t)
	)

	initStateDb := func() *storage.State {
		return &storage.State{
			ID:             stateID,
			IdempotencyKey: createOpts.IdempotencyKey,
			CreatedAt:      createdAt,
			UpdatedAt:      createdAt,
			Status:         uint8(statemachine.NewStatus),
			Step:           string(FirstStep),
			Type:           string(TestType),
			Data: mustBeData(Data{
				Counter: 0,
				Title:   createOpts.Title,
				Amount:  createOpts.Amount,
			}),
		}
	}

	initState := func() *State {
		return &State{
			ID:             stateID,
			IdempotencyKey: createOpts.IdempotencyKey,
			CreatedAt:      createdAt,
			UpdatedAt:      createdAt,
			Status:         statemachine.NewStatus,
			Step:           FirstStep,
			Type:           TestType,
			Data: Data{
				Counter: 0,
				Title:   createOpts.Title,
				Amount:  createOpts.Amount,
			},
		}
	}

	// Создание нового состояния в базе
	t.Run("create new state", func(t *testing.T) {
		deps.uuidGenerator.EXPECT().New().Return(stateID)
		deps.clock.EXPECT().Now().Return(createdAt)
		// По IdempotencyKey ничего не нашли
		deps.storage.EXPECT().GetStateByIdempotencyKey(gomock.Any(), createOpts.IdempotencyKey).Return(nil, storage.ErrNotFound)
		// Создаем запись стейта в базе
		deps.storage.EXPECT().CreateState(gomock.Any(), initStateDb()).Return(nil)

		newState, err := deps.service.Create(deps.ctx, &createOpts)
		require.NoError(t, err)
		require.Equal(t, newState, initState())
	})

	// Попытка создания состояния, но получение уже сохраненного из базы по IdempotencyKey
	t.Run("create return by idempotencyKey", func(t *testing.T) {
		// По IdempotencyKey уже есть запись в базе
		deps.storage.EXPECT().GetStateByIdempotencyKey(gomock.Any(), createOpts.IdempotencyKey).
			Return(initStateDb(), nil)

		newState, err := deps.service.Create(deps.ctx, &createOpts)
		require.NoError(t, err)
		require.Equal(t, newState.ID, stateID)
		require.Equal(t, newState.IdempotencyKey, createOpts.IdempotencyKey)
	})

	// TODO: Complete не существующего стейта

	// Выполняем первый прогон стейт машины
	// Выполняем шаг FirstStep и выходим с ошибкой на шаге TestErrorStep
	t.Run("complete FirstStep and execute error in TestErrorStep", func(t *testing.T) {
		var (
			now = time.Now()
			// Фиксация времени начала выполнения шага
			startExecutedAt1 = now
			startExecutedAt2 = now.Add(2 * time.Second)
			// Фиксация времени выполнения шага
			completeExecutedAt1 = now.Add(time.Second)
			completeExecutedAt2 = now.Add(3 * time.Second)
		)
		// Достаем из базы свежесозданный стейт
		deps.storage.EXPECT().GetStateByID(gomock.Any(), stateID).Return(initStateDb(), nil)

		// Первое выполнение - шаг FirstStep
		deps.clock.EXPECT().Now().Return(startExecutedAt1)
		deps.clock.EXPECT().Now().Return(completeExecutedAt1)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), storage.StepExecuteInfo{
			StateID:            stateID,
			StartExecutedAt:    startExecutedAt1,
			CompleteExecutedAt: completeExecutedAt1,
			Error:              nil,
			PreviewStep:        string(FirstStep),
			NextStep:           lo.ToPtr(string(TestErrorStep)),
		})
		deps.storage.EXPECT().UpdateState(gomock.Any(), storage.UpdateState{
			UpdatedAt: completeExecutedAt1,
			Status:    statemachine.InProgressStatus,
			Step:      string(TestErrorStep),
			Data: mustBeData(Data{
				Counter: 1,
				Title:   "start title",
				Amount:  42,
			}),
		})

		// Второе выполнение - шаг TestErrorStep
		deps.clock.EXPECT().Now().Return(startExecutedAt2)
		deps.clock.EXPECT().Now().Return(completeExecutedAt2)
		//
		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), storage.StepExecuteInfo{
			StateID:            stateID,
			StartExecutedAt:    startExecutedAt2,
			CompleteExecutedAt: completeExecutedAt2,
			Error:              lo.ToPtr("counter eq 2"),
			PreviewStep:        string(TestErrorStep),
		})
		deps.storage.EXPECT().UpdateState(gomock.Any(), storage.UpdateState{
			UpdatedAt: completeExecutedAt1,
			Status:    statemachine.InProgressStatus,
			Step:      string(TestErrorStep),
			Data: mustBeData(Data{
				Counter: 2,
				Title:   "start title",
				Amount:  42,
			}),
		})

		completeState, executeErr, err := deps.service.Complete(deps.ctx, stateID)
		// Должны пройти шаг FirstStep
		// Изменить Counter и Title
		// Далее на шаге TestErrorStep увеличиваем еще счетчик
		// и падаем с ошибкой "counter eq 2"
		require.NoError(t, err)
		require.Error(t, executeErr)
		require.ErrorContains(t, executeErr, "counter eq 2")
		require.Equal(t, completeState.ID, stateID)
		require.Equal(t, completeState.CreatedAt, createdAt)
		// Должна быть равна времени завершения выполнения первого успешного шага те FirstStep
		require.Equal(t, completeState.UpdatedAt, completeExecutedAt1)
		require.Equal(t, completeState.Type, TestType)
		require.Equal(t, completeState.Status, statemachine.InProgressStatus)
		require.Equal(t, completeState.Step, TestErrorStep)
		require.Equal(t, completeState.Data, Data{
			Counter: 2,
			Title:   "start title",
			Amount:  42,
		})
	})

	// Снова выполним TestErrorStep в ожидании что перейдем на Empty
	t.Run("complete TestErrorStep and return Empty", func(t *testing.T) {
		var (
			now = time.Now()
			// Фиксация времени начала выполнения шага
			startExecutedAt = now
			// Фиксация времени выполнения шага
			completeExecutedAt = now.Add(time.Second)
		)
		// Достаем из базы свежесозданный стейт
		init := initStateDb()
		init.Step = string(TestErrorStep)
		init.Status = uint8(statemachine.InProgressStatus)
		init.Data = mustBeData(Data{
			Counter: 2, // Состояние counter как раз что бы получить empty
			Title:   "start title",
			Amount:  42,
		})
		deps.storage.EXPECT().GetStateByID(gomock.Any(), stateID).Return(init, nil)

		// Первое выполнение - шаг TestErrorStep
		deps.clock.EXPECT().Now().Return(startExecutedAt)
		deps.clock.EXPECT().Now().Return(completeExecutedAt)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), storage.StepExecuteInfo{
			StateID:            stateID,
			StartExecutedAt:    startExecutedAt,
			CompleteExecutedAt: completeExecutedAt,
			PreviewStep:        string(TestErrorStep),
		})
		resultData := Data{
			Counter: 3, // Увеличился counter
			Title:   "start title",
			Amount:  42,
		}
		deps.storage.EXPECT().UpdateState(gomock.Any(), storage.UpdateState{
			UpdatedAt: init.UpdatedAt, // Шаг не изменился, по этому время остается старым
			Status:    statemachine.InProgressStatus,
			Step:      string(TestErrorStep),
			Data:      mustBeData(resultData),
		})

		completeState, executeErr, err := deps.service.Complete(deps.ctx, stateID)
		require.NoError(t, err)
		require.NoError(t, executeErr)
		require.Equal(t, completeState.ID, stateID)
		require.Equal(t, completeState.CreatedAt, createdAt)
		// время не изменится
		require.Equal(t, completeState.UpdatedAt, init.UpdatedAt)
		require.Equal(t, completeState.Type, TestType)
		require.Equal(t, completeState.Status, statemachine.InProgressStatus)
		require.Equal(t, completeState.Step, TestErrorStep)
		require.Equal(t, completeState.Data, resultData)
	})

	// Снова выполним TestErrorStep в ожидании перехода на TestNoSaveChangeStep
	// потом на WaitingInputStep
	t.Run("complete TestErrorStep and return Empty", func(t *testing.T) {
		var (
			now = time.Now()
			// Фиксация времени начала выполнения шага
			startExecutedAt = now
			// Фиксация времени выполнения шага
			completeExecutedAt = now.Add(time.Second)
		)
		// Достаем из базы свежесозданный стейт
		init := initStateDb()
		init.Step = string(TestErrorStep)
		init.Status = uint8(statemachine.InProgressStatus)
		init.Data = mustBeData(Data{
			Counter: 3, // Состояние counter как раз что бы перейти на след шаг
			Title:   "start title",
			Amount:  42,
		})
		deps.storage.EXPECT().GetStateByID(gomock.Any(), stateID).Return(init, nil)

		// шаг TestErrorStep
		deps.clock.EXPECT().Now().Return(startExecutedAt)
		deps.clock.EXPECT().Now().Return(completeExecutedAt)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		// Уже нет смысла проверять, тестировали выше
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), gomock.Any())
		deps.storage.EXPECT().UpdateState(gomock.Any(), gomock.Any())

		// Шаг TestNoSaveChangeStep
		// Даже не хотим перетестировать время, оставляем старые значения
		deps.clock.EXPECT().Now().Return(startExecutedAt)
		deps.clock.EXPECT().Now().Return(completeExecutedAt)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		// Уже нет смысла проверять, тестировали выше
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), gomock.Any())
		resultData := Data{
			Counter: 4,             // Увеличился counter (еще в TestErrorStep)
			Title:   "start title", // А значение не должно поменятся
			Amount:  42,
		}
		deps.storage.EXPECT().UpdateState(gomock.Any(), storage.UpdateState{
			UpdatedAt: completeExecutedAt,
			Status:    statemachine.InProgressStatus,
			Step:      string(WaitingInputStep), // Перешли на след шаг
			Data:      mustBeData(resultData),
		}) //

		// шаг WaitingInputStep
		deps.clock.EXPECT().Now().Return(startExecutedAt)
		deps.clock.EXPECT().Now().Return(completeExecutedAt)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		// Уже нет смысла проверять, тестировали выше
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), gomock.Any())
		deps.storage.EXPECT().UpdateState(gomock.Any(), gomock.Any())

		completeState, executeErr, err := deps.service.Complete(deps.ctx, stateID)
		require.NoError(t, err)
		require.NoError(t, executeErr)
		require.Equal(t, completeState.ID, stateID)
		require.Equal(t, completeState.Step, WaitingInputStep)
		require.Equal(t, completeState.Data, resultData)
	})

	// шаг WaitingInputStep, в options пробрасываем левое значение и ждем ошибку выполнения
	t.Run("complete WaitingInputStep with invalid data", func(t *testing.T) {
		var (
			now = time.Now()
			// Фиксация времени начала выполнения шага
			startExecutedAt = now
			// Фиксация времени выполнения шага
			completeExecutedAt = now.Add(time.Second)
		)
		// Достаем из базы свежесозданный стейт
		init := initStateDb()
		init.Step = string(WaitingInputStep)
		init.Status = uint8(statemachine.InProgressStatus)
		init.Data = mustBeData(Data{
			Counter: 3,
			Title:   "start title",
			Amount:  42,
		})
		deps.storage.EXPECT().GetStateByID(gomock.Any(), stateID).Return(init, nil)

		// шаг WaitingInputStep
		deps.clock.EXPECT().Now().Return(startExecutedAt)
		deps.clock.EXPECT().Now().Return(completeExecutedAt)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), gomock.Any())
		deps.storage.EXPECT().UpdateState(gomock.Any(), gomock.Any())

		completeState, executeErr, err := deps.service.Complete(deps.ctx, stateID, "some options")
		require.NoError(t, err)
		require.Error(t, executeErr)
		require.ErrorContains(t, executeErr, "type mismatch")
		require.Equal(t, completeState.ID, stateID)
		require.Equal(t, completeState.Step, WaitingInputStep)
	})

	// шаг WaitingInputStep, в options передаем complete
	t.Run("complete WaitingInputStep with IsComplete", func(t *testing.T) {
		var (
			now = time.Now()
			// Фиксация времени начала выполнения шага
			startExecutedAt = now
			// Фиксация времени выполнения шага
			completeExecutedAt = now.Add(time.Second)
		)
		// Достаем из базы свежесозданный стейт
		init := initStateDb()
		init.Step = string(WaitingInputStep)
		init.Status = uint8(statemachine.InProgressStatus)
		init.Data = mustBeData(Data{
			Counter: 3,
			Title:   "start title",
			Amount:  42,
		})
		deps.storage.EXPECT().GetStateByID(gomock.Any(), stateID).Return(init, nil)

		// шаг WaitingInputStep
		deps.clock.EXPECT().Now().Return(startExecutedAt)
		deps.clock.EXPECT().Now().Return(completeExecutedAt)

		deps.storage.EXPECT().RunTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(context.Context) error) error {
				return txFunc(ctx)
			})
		deps.storage.EXPECT().SaveStepExecuteInfo(gomock.Any(), gomock.Any())
		deps.storage.EXPECT().UpdateState(gomock.Any(), gomock.Any())

		completeState, executeErr, err := deps.service.Complete(deps.ctx, stateID, WaitingInputOptions{
			IsComplete: true,
			NewAmount:  100,
		})
		require.NoError(t, err)
		require.NoError(t, executeErr)
		require.Equal(t, completeState.ID, stateID)
		require.Equal(t, completeState.Status, statemachine.CompletedStatus)
		require.Equal(t, completeState.Step, StepType(""))
		require.Equal(t, completeState.Data, Data{
			Counter: 3,
			Title:   "start title",
			Amount:  100,
		})
	})
}
