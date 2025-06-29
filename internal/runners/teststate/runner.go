package teststate

import (
	"context"
	"github.com/kkiling/torrent2emby/internal/runners"
	"github.com/kkiling/torrent2emby/internal/statemachine"
	"github.com/kkiling/torrent2emby/internal/statemachine/stepper"
)

type CreateState = statemachine.CreateState[Data]
type State = statemachine.State[Data, Status, runners.Type]
type Step = statemachine.Step[Data, Status, runners.Type]
type StepRegistration = statemachine.StepRegistration[Data, Status, runners.Type]
type StepContext = stepper.StepContext[Data, Status, runners.Type]
type StepResult = stepper.StepResult[Data, Status]
type StateMachineService = statemachine.Service[Data, Status, runners.Type, *CreateOptions]

type TaskRunner struct {
}

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{}
}

func (r *TaskRunner) Create(
	_ context.Context,
	options *CreateOptions,
) (CreateState, error) {
	// Логика создания задачи
	data := Data{
		ID:    1,
		Title: options.Title,
	}

	return CreateState{
		Data: data,
	}, nil
}

func (r *TaskRunner) Type() runners.Type {
	return runners.DefaultTaskType
}

func (r *TaskRunner) StepRegistration(_ statemachine.StepRegistrationParams) StepRegistration {
	return StepRegistration{
		Steps: map[Status]Step{
			PendingStatus: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					stepContext.Data.ID = 1
					return stepContext.Empty().WithData(stepContext.Data)
				},
			},
			DoneStatus: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					return stepContext.Empty()
				},
			},
		},
	}
}
