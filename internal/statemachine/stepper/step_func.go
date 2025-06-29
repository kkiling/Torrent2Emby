package stepper

import (
	"context"
	"github.com/kkiling/torrent2emby/internal/statemachine"
)

// stepState управляющее состояние указывающее на результат работы шага
type stepState int

const (
	// emptyStepState шаг не вернул результат своей работы, дальше продвижения не будет, шаг будет выполнен еще раз
	emptyStepState stepState = 0
	// nextStepState шаг вернул результат своей работы, и должен быть продвинут на следующий шаг
	nextStepState stepState = 1
	// failStepState stepState = 2 означает что выпуск нужно зафейлить
	failStepState stepState = 2
	// completeStepState означает что выпуск переведен в статус успешного завершения
	completeStepState stepState = 3
)

// StepContext входные данные функции шага
type StepContext[DataT any, StatusT ~string, TypeT ~string] struct {
	Data            DataT
	State           statemachine.State[DataT, StatusT, TypeT]
	CompleteOptions statemachine.CompleteOptions
}

// Next указываем что нужно перейти на новый статус
func (s *StepContext[DataT, StatusT, TypeT]) Next(status StatusT) *StepResult[DataT, StatusT] {
	return &StepResult[DataT, StatusT]{
		nextStatus: &status,
		state:      nextStepState,
	}
}

// Empty возвращает пустой результат работы шага (стейт машина дальше не продвинется, и шаг будет выполнен еще раз)
func (s *StepContext[DataT, StatusT, TypeT]) Empty() *StepResult[DataT, StatusT] {
	return &StepResult[DataT, StatusT]{
		state: emptyStepState,
	}
}

func (s *StepContext[DataT, StatusT, TypeT]) Error(err error) *StepResult[DataT, StatusT] {
	return &StepResult[DataT, StatusT]{
		state: emptyStepState,
		err:   err,
	}
}

// Fail переводит стейт в терминальное состояние фейла
func (s *StepContext[DataT, StatusT, TypeT]) Fail() *StepResult[DataT, StatusT] {
	return &StepResult[DataT, StatusT]{
		state: failStepState,
	}
}

// Complete переводит стейт в терминальное состояние успеха
func (s *StepContext[DataT, StatusT, TypeT]) Complete() StepResult[DataT, StatusT] {
	return StepResult[DataT, StatusT]{
		state: completeStepState,
	}
}

// StepResult результат работы выпуска
type StepResult[DataT any, StatusT ~string] struct {
	// Указание следующего шага на который должен перейти степпер
	nextStatus *StatusT
	// Новое состояние
	newData *DataT
	// Состояние шага, степер понимает что дальше с ним делать
	state stepState
	// Сохраняем ошибку которая произошла в результате выполнения шага
	err error
}

func (s *StepResult[DataT, StatusT]) WithData(newData DataT) *StepResult[DataT, StatusT] {
	s.newData = &newData
	return s
}

// StepFunc функция выполняющая логику шага
type StepFunc[
	DataT any, StatusT ~string, TypeT ~string,
] func(context.Context, StepContext[DataT, StatusT, TypeT]) *StepResult[DataT, StatusT]
