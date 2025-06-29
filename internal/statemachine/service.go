package statemachine

import (
	"context"
	"github.com/google/uuid"
)

type Config struct {
}

type Service[DataT any, StatusT ~string, TypeT ~string, CreateOptionsT CreateOptions] struct {
	cfg    Config
	runner Runner[DataT, StatusT, TypeT, CreateOptionsT]
}

func NewService[DataT any, StatusT ~string, TypeT ~string, CreateOptionsT CreateOptions](
	cfg Config,
	runner Runner[DataT, StatusT, TypeT, CreateOptionsT],
) *Service[DataT, StatusT, TypeT, CreateOptionsT] {
	return &Service[DataT, StatusT, TypeT, CreateOptionsT]{
		cfg:    cfg,
		runner: runner,
	}
}

// Create создание стейт машины
func (i *Service[DataT, StatusT, TypeT, CreateOptionsT]) Create(
	ctx context.Context,
	options CreateOptionsT,
) (*State[DataT, StatusT, TypeT], error) {

	return nil, nil
}

// Complete выполнение выпуска
func (i *Service[DataT, StatusT, TypeT, CreateOptionsT]) Complete(
	ctx context.Context,
	stateID uuid.UUID,
	options ...CompleteOptions,
) (*State[DataT, StatusT, TypeT], error) {

	return nil, nil
}
