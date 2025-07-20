package teststate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kkiling/torrent2emby/internal/statemachine"
)

type Runner struct {
}

func NewTaskRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Create(
	_ context.Context,
	options *CreateOptions,
) (CreateState, error) {
	// Логика создания задачи
	data := Data{
		Counter: 0,
		Title:   options.Title,
		Amount:  options.Amount,
	}

	return CreateState{
		FirstStep: FirstStep,
		Data:      data,
	}, nil
}

func (r *Runner) Type() Type {
	return TestType
}

func (r *Runner) StepRegistration(_ statemachine.StepRegistrationParams) StepRegistration {
	return StepRegistration{
		Steps: map[StepType]Step{
			FirstStep: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					data.Counter += 1
					data.Title = "start title"
					return stepContext.Next(TestErrorStep).WithData(data)
				},
			},
			TestErrorStep: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					data.Counter += 1
					if data.Counter <= 2 {
						return stepContext.Error(fmt.Errorf("counter eq 2")).WithData(data)
					} else if data.Counter <= 3 {
						return stepContext.Empty().WithData(data)
					} else {
						return stepContext.Next(TestNoSaveChangeStep).WithData(data)
					}
				},
			},
			TestNoSaveChangeStep: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data
					// Изменили значение, но оно не должно записаться в базу
					data.Title = "change title"
					return stepContext.Next(WaitingInputStep)
				},
			},
			WaitingInputStep: {
				OptionsType: reflect.TypeOf(WaitingInputOptions{}), // Устанавливаем тип ожидаемых опций
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение опций выполнения выпуска
					opts := WaitingInputOptions{}
					ok, err := stepContext.GetOptions(&opts)
					if err != nil {
						return stepContext.Error(err)
					}
					if !ok { // Пока не получили опцию, не идем дальше
						return stepContext.Empty()
					}

					data := stepContext.State.Data
					if opts.IsComplete {
						data.Amount = opts.NewAmount
						return stepContext.Complete().WithData(data)
					}

					return stepContext.Fail()
				},
			},
		},
	}
}
