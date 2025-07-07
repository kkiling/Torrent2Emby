package statemachine

import "context"

// Runner интерфейс раннера
type Runner[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string, CreateOptionsT CreateOptions] interface {
	Create(ctx context.Context, options CreateOptionsT) (CreateState[DataT, MetaDataT, StepT], error)
	StepRegistration(params StepRegistrationParams) StepRegistration[DataT, FailDataT, MetaDataT, StepT, TypeT]
	Type() TypeT
}
