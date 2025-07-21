package sqlite

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
	"github.com/kkiling/torrent2emby/internal/statemachine/testutils"
)

func setupTestDB(t *testing.T) *Storage {
	// Инициализируем хранилище
	cfg := Config{
		DSN: testutils.GetSqliteTestDNS(t), // Берем DSN из переменных окружения
	}
	logger := log.NewLogger(log.DebugLevel)

	s, err := NewStorage(cfg, logger)
	require.NoError(t, err)

	// Возвращаем хранилище и функцию очистки
	return s
}

func TestStorage_RunTransaction(t *testing.T) {
	t.Parallel()
	s := setupTestDB(t)
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		stateID := uuid.New()
		// Подготовка тестовых данных
		testState := &storage.State{
			ID:             stateID,
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
		executeInfo := storage.StepExecuteInfo{
			StateID:            stateID,
			StartExecutedAt:    time.Now(),
			CompleteExecutedAt: time.Now(),
			PreviewStep:        "step0",
			NextStep:           lo.ToPtr("step2"),
		}

		// Успешная транзакция
		err := s.RunTransaction(ctx, func(ctxTx context.Context) error {
			// Создаем состояние в транзакции
			if err := s.CreateState(ctxTx, testState); err != nil {
				return err
			}
			return s.SaveStepExecuteInfo(ctxTx, executeInfo)
		})
		require.NoError(t, err)

		// Проверяем, что данные сохранились
		findState, err := s.GetStateByID(ctx, stateID)
		require.NoError(t, err)
		require.Equal(t, findState, testState)

		infos, err := s.GetStepExecuteInfos(ctx, stateID)
		require.NoError(t, err)
		require.Len(t, infos, 1)
		require.Equal(t, infos[0].StateID, executeInfo.StateID)
		require.Equal(t, infos[0].PreviewStep, executeInfo.PreviewStep)
	})

	t.Run("rollback", func(t *testing.T) {
		stateID := uuid.New()
		state := &storage.State{
			ID:             stateID,
			IdempotencyKey: uuid.NewString(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Status:         1,
			Step:           "step1",
			Type:           "type1",
		}

		// Транзакция с ошибкой
		err := s.RunTransaction(ctx, func(ctxTx context.Context) error {
			// Создаем состояние в транзакции
			if err := s.CreateState(ctxTx, state); err != nil {
				return err
			}
			// Намеренно вызываем ошибку
			return errors.New("simulated error")
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "simulated error")

		// Проверяем, что состояние не сохранилось (транзакция откатилась)
		_, err = s.GetStateByID(ctx, stateID)
		require.Error(t, err)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})
}
