package statemachine

import (
	"encoding/json"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

func mapStorageToState[DataT any, StepT ~string, TypeT ~string](state *storage.State) (*State[DataT, StepT, TypeT], error) {
	var data DataT
	err := json.Unmarshal(state.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return &State[DataT, StepT, TypeT]{
		ID:             state.ID,
		IdempotencyKey: state.IdempotencyKey,
		CreatedAt:      state.CreatedAt,
		UpdatedAt:      state.UpdatedAt,
		Status:         Status(state.Status),
		Step:           StepT(state.Step),
		Type:           TypeT(state.Type),
		Data:           data,
	}, nil
}

func mapStateToStorage[DataT any, StepT ~string, TypeT ~string](state *State[DataT, StepT, TypeT]) (*storage.State, error) {
	data, err := json.Marshal(state.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state: %w", err)
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
	}, nil
}
