package teststate

import (
	"github.com/kkiling/torrent2emby/internal/statemachine"
)

type CreateState = statemachine.CreateState[Data, interface{}, StepType]
type State = statemachine.State[Data, interface{}, interface{}, StepType, Type]
type Step = statemachine.Step[Data, interface{}, interface{}, StepType, Type]
type StepRegistration = statemachine.StepRegistration[Data, interface{}, interface{}, StepType, Type]
type StepContext = statemachine.StepContext[Data, interface{}, interface{}, StepType, Type]
type StepResult = statemachine.StepResult[Data, StepType]
type StateMachineService = statemachine.StateMachine[Data, interface{}, interface{}, StepType, Type, *CreateOptions]

func NewState(stateMachineStorage statemachine.Storage) *StateMachineService {
	return statemachine.NewService[Data, interface{}, interface{}, StepType, Type, *CreateOptions](
		statemachine.Config{},
		stateMachineStorage,
		NewTaskRunner(),
	)
}
