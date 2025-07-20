package statemachine

import (
	"context"
	"fmt"
	"reflect"
)

// stepState управляющее состояние указывающее на результат работы шага
type stepState int

const (
	// emptyStepState шаг не вернул результат своей работы, дальше продвижения не будет, шаг будет выполнен еще раз
	emptyStepState stepState = 0
	// errorStepState шаг вернул ошибку в результате своей работы
	errorStepState stepState = 1
	// nextStepState шаг вернул результат своей работы, и должен быть продвинут на следующий шаг
	nextStepState stepState = 2
	// failStepState stepState = 2 означает что выпуск нужно зафейлить
	failStepState stepState = 3
	// completeStepState означает что выпуск переведен в статус успешного завершения
	completeStepState stepState = 4
)

// StepContext входные данные функции шага
type StepContext[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string] struct {
	State               State[DataT, FailDataT, MetaDataT, StepT, TypeT]
	completeOptionsType reflect.Type
	completeOptions     any
}

func (s *StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) GetOptions(v any) (bool, error) {
	// Check if completeOptionsType is not set
	if s.completeOptionsType == nil {
		if s.completeOptions != nil {
			// If completeOptions is set but type is not - return error
			return false, fmt.Errorf("completeOptions is set but completeOptionsType is not defined")
		}
		// Both are not set - return false
		return false, nil
	}

	// Type is set but completeOptions is not
	if s.completeOptions == nil {
		return false, nil
	}

	// Check if the type of completeOptions matches completeOptionsType
	actualType := reflect.TypeOf(s.completeOptions)
	if actualType != s.completeOptionsType {
		return false, fmt.Errorf("type mismatch: expected %v, got %v",
			s.completeOptionsType, actualType)
	}

	// Now perform the conversion
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return false, fmt.Errorf("destination must be a non-nil pointer")
	}

	// Get the value that v points to
	dest := val.Elem()
	src := reflect.ValueOf(s.completeOptions)

	// Check if types are assignable
	if !src.Type().AssignableTo(dest.Type()) {
		return false, fmt.Errorf("cannot assign %v to %v", src.Type(), dest.Type())
	}

	dest.Set(src)
	return true, nil
}

// Next указываем что нужно перейти на новый статус
func (s *StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) Next(status StepT) *StepResult[DataT, StepT] {
	return &StepResult[DataT, StepT]{
		nextStatus: &status,
		state:      nextStepState,
	}
}

// Empty возвращает пустой результат работы шага (стейт машина дальше не продвинется, и шаг будет выполнен еще раз)
func (s *StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) Empty() *StepResult[DataT, StepT] {
	return &StepResult[DataT, StepT]{
		state: emptyStepState,
	}
}

func (s *StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) Error(err error) *StepResult[DataT, StepT] {
	return &StepResult[DataT, StepT]{
		state: errorStepState,
		err:   err,
	}
}

// Fail переводит стейт в терминальное состояние фейла
func (s *StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) Fail() *StepResult[DataT, StepT] {
	return &StepResult[DataT, StepT]{
		state: failStepState,
	}
}

// Complete переводит стейт в терминальное состояние успеха
func (s *StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) Complete() *StepResult[DataT, StepT] {
	return &StepResult[DataT, StepT]{
		state: completeStepState,
	}
}

// StepResult результат работы выпуска
type StepResult[DataT any, StepT ~string] struct {
	// Указание следующего шага на который должен перейти степпер
	nextStatus *StepT
	// Новое состояние
	newData *DataT
	// Состояние шага, степер понимает что дальше с ним делать
	state stepState
	// Сохраняем ошибку которая произошла в результате выполнения шага
	err error
}

func (s *StepResult[DataT, StepT]) WithData(newData DataT) *StepResult[DataT, StepT] {
	s.newData = &newData
	return s
}

// StepFunc функция выполняющая логику шага
type StepFunc[
	DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string,
] func(context.Context, StepContext[DataT, FailDataT, MetaDataT, StepT, TypeT]) *StepResult[DataT, StepT]
