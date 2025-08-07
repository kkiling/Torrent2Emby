package tvshowdeliverystate

import (
	"github.com/kkiling/statemachine"

	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners"
)

type CreateState = statemachine.CreateState[TVShowDeliveryData, runners.Metadata, StepDelivery]
type State = statemachine.State[TVShowDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type Step = statemachine.Step[TVShowDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type StepRegistration = statemachine.StepRegistration[TVShowDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type StepContext = statemachine.StepContext[TVShowDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type]
type StepResult = statemachine.StepResult[TVShowDeliveryData, StepDelivery]
type StateMachineService = statemachine.StateMachine[TVShowDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type, CreateOptions]

func NewState(contentDelivery ContentDelivery, stateMachineStorage statemachine.Storage) *StateMachineService {
	return statemachine.NewService[TVShowDeliveryData, runners.FailData, runners.Metadata, StepDelivery, runners.Type, CreateOptions](
		statemachine.Config{},
		stateMachineStorage,
		NewTaskRunner(contentDelivery),
	)
}
