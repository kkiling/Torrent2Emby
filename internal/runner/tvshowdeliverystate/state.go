package tvshowdeliverystate

import (
	"github.com/kkiling/torrent2emby/internal/statemachine"
)

type CreateState = statemachine.CreateState[ContentDeliveryData, ContentDeliveryMetadata, StepDelivery]
type State = statemachine.State[ContentDeliveryData, ContentDeliveryFailData, ContentDeliveryMetadata, StepDelivery, Type]
type Step = statemachine.Step[ContentDeliveryData, ContentDeliveryFailData, ContentDeliveryMetadata, StepDelivery, Type]
type StepRegistration = statemachine.StepRegistration[ContentDeliveryData, ContentDeliveryFailData, ContentDeliveryMetadata, StepDelivery, Type]
type StepContext = statemachine.StepContext[ContentDeliveryData, ContentDeliveryFailData, ContentDeliveryMetadata, StepDelivery, Type]
type StepResult = statemachine.StepResult[ContentDeliveryData, StepDelivery]
type StateMachineService = statemachine.StateMachine[ContentDeliveryData, ContentDeliveryFailData, ContentDeliveryMetadata, StepDelivery, Type, CreateOptions]

func NewState(contentDelivery ContentDelivery, stateMachineStorage statemachine.Storage) *StateMachineService {
	return statemachine.NewService[ContentDeliveryData, ContentDeliveryFailData, ContentDeliveryMetadata, StepDelivery, Type, CreateOptions](
		statemachine.Config{},
		stateMachineStorage,
		NewTaskRunner(contentDelivery),
	)
}
