package teststate

import (
	"context"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage/sqlitestorage"
	"github.com/kkiling/torrent2emby/internal/statemachine/testutils"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/golang/mock/gomock"

	mock_statemachine "github.com/kkiling/torrent2emby/internal/statemachine/mocks"
)

type testDeps struct {
	ctx           context.Context
	ctrl          *gomock.Controller
	service       *StateMachineService
	clock         *mock_statemachine.MockClock
	uuidGenerator *mock_statemachine.MockUUIDGenerator
	storageMock   *mock_statemachine.MockStorage
	storageSqlite *sqlitestorage.Storage
}

func setupTestDeps(t *testing.T, opts ...func(d *testDeps)) *testDeps {
	deps := &testDeps{
		ctx:  context.Background(),
		ctrl: gomock.NewController(t),
	}

	if deps.storageMock == nil {
		deps.storageMock = mock_statemachine.NewMockStorage(deps.ctrl)
	}
	if deps.uuidGenerator == nil {
		deps.uuidGenerator = mock_statemachine.NewMockUUIDGenerator(deps.ctrl)
	}
	if deps.clock == nil {
		deps.clock = mock_statemachine.NewMockClock(deps.ctrl)
	}

	deps.service = NewState(deps.storageMock)
	deps.service.SetClock(deps.clock)
	deps.service.SetUUIDGenerator(deps.uuidGenerator)

	for _, opt := range opts {
		opt(deps)
	}

	return deps
}

func setupTestDB(t *testing.T) *sqlitestorage.Storage {
	// Инициализируем хранилище
	cfg := sqlitestorage.Config{
		DSN: testutils.GetSqliteTestDNS(t), // Берем DSN из переменных окружения
	}
	logger := log.NewLogger(log.DebugLevel)

	s, err := sqlitestorage.NewStorage(cfg, logger)
	require.NoError(t, err)

	// Возвращаем хранилище и функцию очистки
	return s
}

func setupTestDepsSqlite(t *testing.T) *testDeps {
	return setupTestDeps(t, func(deps *testDeps) {
		deps.storageSqlite = setupTestDB(t)
		deps.storageMock = nil
		deps.service = NewState(deps.storageSqlite)
		deps.service.SetClock(deps.clock)
		deps.service.SetUUIDGenerator(deps.uuidGenerator)
	})
}
