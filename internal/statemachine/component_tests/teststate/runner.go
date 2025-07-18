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
					stepContext.Data.Counter += 1
					stepContext.Data.Title = "start title"
					return stepContext.Next(TestErrorStep).WithData(stepContext.Data)
				},
			},
			TestErrorStep: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					stepContext.Data.Counter += 1
					if stepContext.Data.Counter <= 2 {
						return stepContext.Error(fmt.Errorf("counter eq 2")).WithData(stepContext.Data)
					} else if stepContext.Data.Counter <= 3 {
						return stepContext.Empty().WithData(stepContext.Data)
					} else {
						return stepContext.Next(TestNoSaveChangeStep).WithData(stepContext.Data)
					}
				},
			},
			TestNoSaveChangeStep: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Изменили значение, но оно не должно записаться в базу
					stepContext.Data.Title = "change title"
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

					if opts.IsComplete {
						stepContext.Data.Amount = opts.NewAmount
						return stepContext.Complete().WithData(stepContext.Data)
					}

					return stepContext.Fail()
				},
			},
		},
	}
}
