package sqlite

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

func (s *Storage) CreateState(ctx context.Context, state *storage.State) error {
	query := `
		INSERT INTO state (
			id, idempotency_key, created_at, updated_at, 
			status, step, type, data, fail_data, meta_data
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.next(ctx).ExecContext(ctx, query,
		state.ID[:],
		state.IdempotencyKey[:],
		state.CreatedAt,
		state.UpdatedAt,
		state.Status,
		state.Step,
		state.Type,
		state.Data,
		state.FailData,
		state.MetaData,
	)

	if err != nil {
		return handleError(err)
	}

	return err
}

func (s *Storage) GetStateByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*storage.State, error) {
	query := `
		SELECT 
			id, idempotency_key, created_at, updated_at, 
			status, step, type, data, fail_data, meta_data
		FROM state
		WHERE idempotency_key = ?
	`

	var state storage.State
	var idBytes, keyBytes []byte

	err := s.next(ctx).QueryRowContext(ctx, query, idempotencyKey[:]).Scan(
		&idBytes,
		&keyBytes,
		&state.CreatedAt,
		&state.UpdatedAt,
		&state.Status,
		&state.Step,
		&state.Type,
		&state.Data,
		&state.FailData,
		&state.MetaData,
	)

	if err != nil {
		return nil, handleError(err)
	}

	state.ID, err = uuid.FromBytes(idBytes)
	if err != nil {
		return nil, err
	}

	state.IdempotencyKey, err = uuid.FromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (s *Storage) GetStateByID(ctx context.Context, stateID uuid.UUID) (*storage.State, error) {
	query := `
		SELECT 
			id, idempotency_key, created_at, updated_at, 
			status, step, type, data, fail_data, meta_data
		FROM state
		WHERE id = ?
	`

	var state storage.State
	var idBytes, keyBytes []byte

	err := s.next(ctx).QueryRowContext(ctx, query, stateID[:]).Scan(
		&idBytes,
		&keyBytes,
		&state.CreatedAt,
		&state.UpdatedAt,
		&state.Status,
		&state.Step,
		&state.Type,
		&state.Data,
		&state.FailData,
		&state.MetaData,
	)

	if err != nil {
		return nil, handleError(err)
	}

	state.ID, err = uuid.FromBytes(idBytes)
	if err != nil {
		return nil, err
	}

	state.IdempotencyKey, err = uuid.FromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (s *Storage) SaveStepExecuteInfo(ctx context.Context, execute storage.StepExecuteInfo) error {
	query := `
		INSERT INTO step_execute_info (
			state_id, start_executed_at, complete_executed_at, 
			error, preview_step, next_step
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.next(ctx).ExecContext(ctx, query,
		execute.StateID[:],
		execute.StartExecutedAt,
		execute.CompleteExecutedAt,
		execute.Error,
		execute.PreviewStep,
		execute.NextStep,
	)

	return handleError(err)
}

func (s *Storage) GetStepExecuteInfos(ctx context.Context, stateID uuid.UUID) ([]storage.StepExecuteInfo, error) {
	query := `
        SELECT 
            state_id, start_executed_at, complete_executed_at,
            error, preview_step, next_step
        FROM step_execute_info
        WHERE state_id = ?
        ORDER BY start_executed_at
    `

	rows, err := s.next(ctx).QueryContext(ctx, query, stateID[:])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var infos []storage.StepExecuteInfo
	for rows.Next() {
		var info storage.StepExecuteInfo
		var stateIDBytes []byte
		var errorStr sql.NullString
		var nextStep sql.NullString
		var completeTime sql.NullTime

		err := rows.Scan(
			&stateIDBytes,
			&info.StartExecutedAt,
			&completeTime,
			&errorStr,
			&info.PreviewStep,
			&nextStep,
		)
		if err != nil {
			return nil, handleError(err)
		}

		info.StateID, err = uuid.FromBytes(stateIDBytes)
		if err != nil {
			return nil, err
		}

		if completeTime.Valid {
			info.CompleteExecutedAt = completeTime.Time
		}

		if errorStr.Valid {
			info.Error = &errorStr.String
		}

		if nextStep.Valid {
			info.NextStep = &nextStep.String
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func (s *Storage) UpdateState(ctx context.Context, stateID uuid.UUID, state storage.UpdateState) error {
	query := `
		UPDATE state
		SET 
			updated_at = ?,
			status = ?,
			step = ?,
			data = ?,
			fail_data = ?,
			meta_data = ?
		WHERE id = ?
	`

	result, err := s.next(ctx).ExecContext(ctx, query,
		state.UpdatedAt,
		state.Status,
		state.Step,
		state.Data,
		state.FailData,
		state.MetaData,
		stateID[:],
	)

	if err != nil {
		return handleError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return handleError(err)
	}

	if rowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}
