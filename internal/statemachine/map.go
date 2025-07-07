package statemachine

import (
	"encoding/json"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

func mapStorageToState[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string](
	state *storage.State,
) (*State[DataT, FailDataT, MetaDataT, StepT, TypeT], error) {
	var err error
	var data DataT
	if len(state.Data) > 0 {
		err := json.Unmarshal(state.Data, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	var failData FailDataT
	if len(state.FailData) > 0 {
		err = json.Unmarshal(state.FailData, &failData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal failData: %w", err)
		}
	}

	var metaData MetaDataT

	if len(state.MetaData) > 0 {
		err = json.Unmarshal(state.MetaData, &metaData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metaData: %w", err)
		}
	}

	return &State[DataT, FailDataT, MetaDataT, StepT, TypeT]{
		ID:             state.ID,
		IdempotencyKey: state.IdempotencyKey,
		CreatedAt:      state.CreatedAt,
		UpdatedAt:      state.UpdatedAt,
		Status:         Status(state.Status),
		Step:           StepT(state.Step),
		Type:           TypeT(state.Type),
		Data:           data,
		FailData:       failData,
		MetaData:       metaData,
	}, nil
}

func mapStateToStorage[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string](
	state *State[DataT, FailDataT, MetaDataT, StepT, TypeT],
) (*storage.State, error) {

	var err error
	var data []byte

	data, err = json.Marshal(state.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	if string(data) == "null" {
		data = []byte{}
	}

	var failData []byte
	failData, err = json.Marshal(state.FailData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal failData: %w", err)
	}
	if string(failData) == "null" {
		failData = []byte{}
	}

	var metaData []byte
	metaData, err = json.Marshal(state.MetaData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metaData: %w", err)
	}
	if string(metaData) == "null" {
		metaData = []byte{}
	}

	return &storage.State{
		ID:             state.ID,
		IdempotencyKey: state.IdempotencyKey,
		CreatedAt:      state.CreatedAt,
		UpdatedAt:      state.UpdatedAt,
		Status:         uint8(state.Status),
		Step:           string(state.Step),
		Type:           string(state.Type),
		Data:           data,
		FailData:       failData,
		MetaData:       metaData,
	}, nil
}
