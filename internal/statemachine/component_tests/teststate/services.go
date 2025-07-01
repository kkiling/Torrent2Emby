package teststate

import (
	"github.com/kkiling/torrent2emby/internal/statemachine"
)

type CreateState = statemachine.CreateState[Data, StepType]
type State = statemachine.State[Data, StepType, Type]
type Step = statemachine.Step[Data, StepType, Type]
type StepRegistration = statemachine.StepRegistration[Data, StepType, Type]
type StepContext = statemachine.StepContext[Data, StepType, Type]
type StepResult = statemachine.StepResult[Data, StepType]
type StateMachineService = statemachine.StateMachine[Data, StepType, Type, *CreateOptions]

func NewState(stateMachineStorage statemachine.Storage) *StateMachineService {
	return statemachine.NewService[Data, StepType, Type, *CreateOptions](
		statemachine.Config{},
		stateMachineStorage,
		NewTaskRunner(),
	)
}
