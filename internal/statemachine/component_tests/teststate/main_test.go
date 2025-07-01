package teststate

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_statemachine "github.com/kkiling/torrent2emby/internal/statemachine/mocks"
	"testing"
)

type testDeps struct {
	ctx           context.Context
	ctrl          *gomock.Controller
	service       *StateMachineService
	clock         *mock_statemachine.MockClock
	uuidGenerator *mock_statemachine.MockUUIDGenerator
	storage       *mock_statemachine.MockStorage
}

func setupTestDeps(t *testing.T, _ ...func(d *testDeps)) *testDeps {
	deps := &testDeps{
		ctx:  context.Background(),
		ctrl: gomock.NewController(t),
	}

	if deps.storage == nil {
		deps.storage = mock_statemachine.NewMockStorage(deps.ctrl)
	}
	if deps.uuidGenerator == nil {
		deps.uuidGenerator = mock_statemachine.NewMockUUIDGenerator(deps.ctrl)
	}
	if deps.clock == nil {
		deps.clock = mock_statemachine.NewMockClock(deps.ctrl)
	}

	deps.service = NewState(deps.storage)
	deps.service.SetClock(deps.clock)
	deps.service.SetUUIDGenerator(deps.uuidGenerator)

	return deps
}
