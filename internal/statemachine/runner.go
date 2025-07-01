package statemachine

import "context"

// Runner интерфейс раннера
type Runner[DataT any, StepT ~string, TypeT ~string, CreateOptionsT CreateOptions] interface {
	Create(ctx context.Context, options CreateOptionsT) (CreateState[DataT, StepT], error)
	StepRegistration(params StepRegistrationParams) StepRegistration[DataT, StepT, TypeT]
	Type() TypeT
}
