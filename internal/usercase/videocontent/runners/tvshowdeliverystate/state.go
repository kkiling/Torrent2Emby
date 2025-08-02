package tvshowdeliverystate

import (
	"github.com/kkiling/statemachine"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners"
)

type CreateState = statemachine.CreateState[ContentDeliveryData, runners.Metadata, StepDelivery]
type State = statemachine.State[ContentDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type Step = statemachine.Step[ContentDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type StepRegistration = statemachine.StepRegistration[ContentDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type StepContext = statemachine.StepContext[ContentDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type StepResult = statemachine.StepResult[ContentDeliveryData, StepDelivery]
type StateMachineService = statemachine.StateMachine[ContentDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type, CreateOptions]

func NewState(contentDelivery ContentDelivery, stateMachineStorage statemachine.Storage) *StateMachineService {
	return statemachine.NewService[ContentDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type, CreateOptions](
		statemachine.Config{},
		stateMachineStorage,
		NewTaskRunner(contentDelivery),
	)
}
