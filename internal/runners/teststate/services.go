package teststate

import (
	"github.com/kkiling/torrent2emby/internal/runners"
	"github.com/kkiling/torrent2emby/internal/statemachine"
)

func NewState() *StateMachineService {
	return statemachine.NewService[Data, Status, runners.Type, *CreateOptions](
		statemachine.Config{},
		NewTaskRunner(),
	)
}
