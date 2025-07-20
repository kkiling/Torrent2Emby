package sqlite

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

func TestCreateState(t *testing.T) {
	t.Parallel()

	s := setupTestDB(t)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()

		state := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         1,
			Step:           "initial",
			Type:           "test",
			Data:           []byte(`{"key":"value"}`),
			FailData:       nil,
			MetaData:       []byte(`{"meta":"data"}`),
		}

		err := s.CreateState(ctx, state)
		require.NoError(t, err)

		// Проверяем, что состояние действительно сохранено
		savedState, err := s.GetStateByID(ctx, state.ID)
		require.NoError(t, err)
		require.Equal(t, state, savedState)
	})

	t.Run("duplicate id", func(t *testing.T) {
		t.Parallel()

		state := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         1,
			Step:           "initial",
			Type:           "test",
			Data:           []byte(`{"key":"value"}`),
			FailData:       nil,
			MetaData:       []byte(`{"meta":"data"}`),
		}

		// Первое сохранение должно быть успешным
		err := s.CreateState(ctx, state)
		require.NoError(t, err)

		// Попытка сохранить с тем же ID должна вернуть ошибку
		err = s.CreateState(ctx, state)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrAlreadyExists)
	})

	t.Run("duplicate idempotency_key", func(t *testing.T) {
		t.Parallel()
		idempotencyKey := uuid.NewString()

		state1 := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: idempotencyKey,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         1,
			Step:           "initial",
			Type:           "test",
			Data:           []byte(`{"key":"value1"}`),
			FailData:       nil,
			MetaData:       []byte(`{"meta":"data1"}`),
		}

		state2 := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: idempotencyKey,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         2,
			Step:           "next",
			Type:           "test",
			Data:           []byte(`{"key":"value2"}`),
			FailData:       nil,
			MetaData:       []byte(`{"meta":"data2"}`),
		}

		// Первое сохранение должно быть успешным
		err := s.CreateState(ctx, state1)
		require.NoError(t, err)

		// Попытка сохранить с тем же ключом идемпотентности должна вернуть ошибку
		err = s.CreateState(ctx, state2)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrAlreadyExists)
	})

	t.Run("with empty data", func(t *testing.T) {
		t.Parallel()

		state := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         1,
			Step:           "empty",
			Type:           "test",
			Data:           nil,
			FailData:       nil,
			MetaData:       nil,
		}

		err := s.CreateState(ctx, state)
		require.NoError(t, err)

		// Проверяем, что состояние с пустыми данными сохранилось
		savedState, err := s.GetStateByID(ctx, state.ID)
		require.NoError(t, err)
		require.Nil(t, savedState.Data)
		require.Nil(t, savedState.FailData)
		require.Nil(t, savedState.MetaData)
	})
}

func TestGetStateByIdempotencyKey(t *testing.T) {
	t.Parallel()
	s := setupTestDB(t)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		// Подготовка тестовых данных
		testState := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC().Truncate(time.Second), // Округляем до секунд для сравнения
			UpdatedAt:      time.Now().UTC().Truncate(time.Second),
			Status:         1,
			Step:           "initial",
			Type:           "test",
			Data:           []byte(`{"key":"value"}`),
			FailData:       nil,
			MetaData:       []byte(`{"meta":"data"}`),
		}

		// Сначала создаем состояние
		err := s.CreateState(ctx, testState)
		require.NoError(t, err)

		// Получаем состояние по ключу идемпотентности
		state, err := s.GetStateByIdempotencyKey(ctx, testState.IdempotencyKey)
		require.NoError(t, err)
		require.NotNil(t, state)
		require.Equal(t, testState, state)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// Пытаемся получить несуществующее состояние
		nonExistentKey := uuid.NewString()
		state, err := s.GetStateByIdempotencyKey(ctx, nonExistentKey)

		// Проверяем, что получили ожидаемую ошибку
		require.Nil(t, state)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})

	t.Run("empty idempotency key", func(t *testing.T) {
		t.Parallel()

		// Пытаемся получить состояние с пустым UUID
		state, err := s.GetStateByIdempotencyKey(ctx, "")

		require.Nil(t, state)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})
}

func TestGetStateByID(t *testing.T) {
	t.Parallel()
	s := setupTestDB(t)
	ctx := context.Background()

	// Подготовка тестовых данных
	testState := &storage.State{
		ID:             uuid.New(),
		IdempotencyKey: uuid.NewString(),
		CreatedAt:      time.Now().UTC().Truncate(time.Second), // Округляем для точного сравнения
		UpdatedAt:      time.Now().UTC().Truncate(time.Second),
		Status:         124,
		Step:           "initial",
		Type:           "test_transaction",
		Data:           []byte(`{"amount":100}`),
		FailData:       nil,
		MetaData:       []byte(`{"user_id":123}`),
	}

	// Сначала создаем состояние
	err := s.CreateState(ctx, testState)
	require.NoError(t, err)

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()
		// Получаем состояние по ID
		state, err := s.GetStateByID(ctx, testState.ID)
		require.NoError(t, err)
		require.NotNil(t, state)

		// Проверяем все поля
		require.Equal(t, testState, state)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		// Пытаемся получить несуществующее состояние
		nonExistentID := uuid.New()
		state, err := s.GetStateByID(ctx, nonExistentID)

		require.Nil(t, state)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})

	t.Run("empty UUID", func(t *testing.T) {
		t.Parallel()
		// Пытаемся получить состояние с пустым UUID
		state, err := s.GetStateByID(ctx, uuid.Nil)

		require.Nil(t, state)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})
}

func TestSaveStepExecuteInfo(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	saveTestState := func(t *testing.T) *storage.State {
		// Создаем тестовое состояние
		testState := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         1,
			Step:           "initial",
			Type:           "test",
		}
		require.NoError(t, s.CreateState(ctx, testState))

		return testState
	}

	t.Run("successful save", func(t *testing.T) {
		testState := saveTestState(t)

		info := storage.StepExecuteInfo{
			StateID:            testState.ID,
			StartExecutedAt:    time.Now().UTC().Truncate(time.Second),
			CompleteExecutedAt: time.Now().UTC().Add(5 * time.Second).Truncate(time.Second),
			Error:              lo.ToPtr("test error"),
			PreviewStep:        "prev_step",
			NextStep:           lo.ToPtr("next_step"),
		}

		err := s.SaveStepExecuteInfo(ctx, info)
		require.NoError(t, err)

		// Проверяем через GetStepExecuteInfos
		saved, err := s.GetStepExecuteInfos(ctx, testState.ID)
		require.NoError(t, err)
		require.Len(t, saved, 1)
		require.Equal(t, info, saved[0])
	})

	t.Run("fail on non-existent state", func(t *testing.T) {
		info := storage.StepExecuteInfo{
			StateID:         uuid.New(), // Несуществующий ID
			StartExecutedAt: time.Now(),
			PreviewStep:     "test",
		}

		err := s.SaveStepExecuteInfo(ctx, info)
		require.Error(t, err)
		require.Contains(t, err.Error(), "FOREIGN KEY constraint failed")
	})

	t.Run("successful save multiple steps", func(t *testing.T) {
		testState := saveTestState(t)

		info := storage.StepExecuteInfo{
			StateID:            testState.ID,
			StartExecutedAt:    time.Now().UTC().Truncate(time.Second),
			CompleteExecutedAt: time.Now().UTC().Add(5 * time.Second).Truncate(time.Second),
			PreviewStep:        "1 step",
			NextStep:           lo.ToPtr("2 step"),
		}

		err := s.SaveStepExecuteInfo(ctx, info)
		require.NoError(t, err)

		info.PreviewStep = "2 step"
		info.NextStep = nil
		err = s.SaveStepExecuteInfo(ctx, info)
		require.NoError(t, err)

		info.PreviewStep = "2 step"
		info.NextStep = lo.ToPtr("3 step")
		err = s.SaveStepExecuteInfo(ctx, info)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrAlreadyExists) // Пытаемся сохранить стейт с тем же шагом PreviewStep

		// Проверяем через GetStepExecuteInfos
		saved, err := s.GetStepExecuteInfos(ctx, testState.ID)
		require.NoError(t, err)
		require.Len(t, saved, 2)
		require.Equal(t, saved[0].PreviewStep, "1 step")
		require.Equal(t, saved[1].PreviewStep, "2 step")
	})
}

func TestGetStepExecuteInfos(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	saveTestState := func(t *testing.T) *storage.State {
		// Создаем тестовое состояние
		testState := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			Status:         1,
			Step:           "initial",
			Type:           "test",
		}
		require.NoError(t, s.CreateState(ctx, testState))

		return testState
	}

	t.Run("get single step info", func(t *testing.T) {
		state := saveTestState(t)
		info := storage.StepExecuteInfo{
			StateID:            state.ID,
			StartExecutedAt:    time.Now().UTC().Truncate(time.Second),
			CompleteExecutedAt: time.Now().UTC().Add(5 * time.Second).Truncate(time.Second),
			Error:              lo.ToPtr("test error"),
			PreviewStep:        "initial",
			NextStep:           lo.ToPtr("next_step"),
		}

		err := s.SaveStepExecuteInfo(ctx, info)
		require.NoError(t, err)

		infos, err := s.GetStepExecuteInfos(ctx, state.ID)
		require.NoError(t, err)
		require.Len(t, infos, 1)

		require.Equal(t, info, infos[0])
	})

	t.Run("get multiple steps ordered by time", func(t *testing.T) {
		state := saveTestState(t)

		now := time.Now().UTC().Truncate(time.Second)
		infos := []storage.StepExecuteInfo{
			{
				StateID:            state.ID,
				StartExecutedAt:    now.Add(10 * time.Second),
				CompleteExecutedAt: now.Add(15 * time.Second),
				PreviewStep:        "step_2",
				NextStep:           lo.ToPtr("step_3"),
			},
			{
				StateID:            state.ID,
				StartExecutedAt:    now,
				CompleteExecutedAt: now.Add(5 * time.Second),
				PreviewStep:        "step_1",
				NextStep:           lo.ToPtr("step_2"),
			},
			{
				StateID:            state.ID,
				StartExecutedAt:    now.Add(20 * time.Second),
				CompleteExecutedAt: now.Add(25 * time.Second),
				PreviewStep:        "step_3",
				NextStep:           nil,
			},
		}

		for _, info := range infos {
			err := s.SaveStepExecuteInfo(ctx, info)
			require.NoError(t, err)
		}

		find, err := s.GetStepExecuteInfos(ctx, state.ID)
		require.NoError(t, err)
		require.Len(t, infos, 3)
		require.Equal(t, infos[1], find[0]) // step_1
		require.Equal(t, infos[0], find[1]) // step_2
		require.Equal(t, infos[2], find[2]) // step_3
	})

	t.Run("return empty slice for unknown state", func(t *testing.T) {
		infos, err := s.GetStepExecuteInfos(ctx, uuid.New())
		require.NoError(t, err)
		require.Empty(t, infos)
	})
}

func TestUpdateState(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	const (
		StateStatusPending    = 12
		StateStatusCompleted  = 3
		StateStatusFailed     = 7
		StateStatusProcessing = 4
	)

	// Вспомогательная функция для создания тестового состояния
	createTestState := func(t *testing.T) *storage.State {
		state := &storage.State{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now().UTC().Truncate(time.Second),
			UpdatedAt:      time.Now().UTC().Truncate(time.Second),
			Status:         StateStatusPending,
			Step:           "initial",
			Type:           "test",
			Data:           []byte(`{"initial": true}`),
		}
		require.NoError(t, s.CreateState(ctx, state))
		return state
	}

	t.Run("successful update all fields", func(t *testing.T) {
		testState := createTestState(t)
		newTime := time.Now().UTC().Truncate(time.Second)

		update := storage.UpdateState{
			UpdatedAt: newTime,
			Status:    StateStatusCompleted,
			Step:      "completed",
			Data:      []byte(`{"completed": true}`),
			FailData:  nil,
			MetaData:  []byte(`{"meta": "data"}`),
		}

		err := s.UpdateState(ctx, testState.ID, update)
		require.NoError(t, err)

		// Проверяем обновленные данные
		updatedState, err := s.GetStateByID(ctx, testState.ID)
		require.NoError(t, err)

		require.Equal(t, update.UpdatedAt, updatedState.UpdatedAt)
		require.Equal(t, update.Status, updatedState.Status)
		require.Equal(t, update.Step, updatedState.Step)
		require.Equal(t, update.Data, updatedState.Data)
		require.Equal(t, update.FailData, updatedState.FailData)
		require.Equal(t, update.MetaData, updatedState.MetaData)
	})

	t.Run("successful partial update", func(t *testing.T) {
		testState := createTestState(t)
		newTime := time.Now().UTC().Truncate(time.Second)

		update := storage.UpdateState{
			UpdatedAt: newTime,
			Status:    StateStatusFailed,
			Step:      "failed",
			FailData:  []byte(`{"error": "something went wrong"}`),
			// Data и MetaDataT не обновляем
			Data:     testState.Data,
			MetaData: testState.MetaData,
		}

		err := s.UpdateState(ctx, testState.ID, update)
		require.NoError(t, err)

		// Проверяем обновленные данные
		updatedState, err := s.GetStateByID(ctx, testState.ID)
		require.NoError(t, err)

		require.Equal(t, update.UpdatedAt, updatedState.UpdatedAt)
		require.Equal(t, update.Status, updatedState.Status)
		require.Equal(t, update.Step, updatedState.Step)
		require.Equal(t, update.FailData, updatedState.FailData)
		// Проверяем, что остальные поля не изменились
		require.Equal(t, testState.Data, updatedState.Data)
		require.Equal(t, testState.MetaData, updatedState.MetaData)
	})

	t.Run("fail on non-existent state", func(t *testing.T) {
		update := storage.UpdateState{
			UpdatedAt: time.Now(),
			Status:    StateStatusPending,
			Step:      "test",
		}

		err := s.UpdateState(ctx, uuid.New(), update)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})

	t.Run("update with null data fields", func(t *testing.T) {
		testState := createTestState(t)
		// Предварительно заполняем данные
		require.NoError(t, s.UpdateState(ctx, testState.ID, storage.UpdateState{
			UpdatedAt: time.Now(),
			Status:    StateStatusPending,
			Step:      "has_data",
			Data:      []byte(`{"test":1}`),
			FailData:  []byte(`{"error":"test"}`),
			MetaData:  []byte(`{"meta":1}`),
		}))

		// Обновляем с null-значениями
		update := storage.UpdateState{
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
			Status:    StateStatusProcessing,
			Step:      "null_data",
			Data:      nil,
			FailData:  nil,
			MetaData:  nil,
		}

		err := s.UpdateState(ctx, testState.ID, update)
		require.NoError(t, err)

		// Проверяем, что данные обнулились
		updatedState, err := s.GetStateByID(ctx, testState.ID)
		require.NoError(t, err)

		require.Nil(t, updatedState.Data)
		require.Nil(t, updatedState.FailData)
		require.Nil(t, updatedState.MetaData)
	})

	t.Run("context cancellation", func(t *testing.T) {
		testState := createTestState(t)
		ctx, cancel := context.WithCancel(ctx)
		cancel() // сразу отменяем контекст

		update := storage.UpdateState{
			UpdatedAt: time.Now(),
			Status:    StateStatusPending,
			Step:      "canceled",
		}

		err := s.UpdateState(ctx, testState.ID, update)
		require.Error(t, err)
		require.True(t, errors.Is(err, context.Canceled))
	})
}
