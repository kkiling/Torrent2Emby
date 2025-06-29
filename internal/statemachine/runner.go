package statemachine

import "context"

// Runner интерфейс раннера
type Runner[DataT any, StatusT ~string, TypeT ~string, CreateOptionsT CreateOptions] interface {

	//Complete(ctx context.Context,
	//	state *State[DataT, StatusT, TypeT],
	//	options ...CompleteOptions,
	//) (*State[DataT, StatusT, TypeT], error)

	Create(ctx context.Context, options CreateOptionsT) (CreateState[DataT], error)
	Type() TypeT
	StepRegistration(params StepRegistrationParams) StepRegistration[DataT, StatusT, TypeT]
}
