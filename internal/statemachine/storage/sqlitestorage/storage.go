package sqlitestorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type txKey struct{}

type Config struct {
	DSN string
}

type Storage struct {
	config Config
	db     *sql.DB
	logger log.Logger
}

func NewStorage(config Config, logger log.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Включаем поддержку внешних ключей
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed PRAGMA foreign_keys = ON;: %w", err)
	}

	return &Storage{config: config, db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNotFound
	}

	// Проверяем, является ли ошибка ошибкой уникальности
	sqliteErr := sqlite3.Error{}
	if ok := errors.As(err, &sqliteErr); ok {
		if strings.Contains(sqliteErr.Error(), "FOREIGN KEY constraint failed") {
			return err
		}
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return storage.ErrAlreadyExists
		}
	}

	return err
}

func (s *Storage) next(ctx context.Context) SQLExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		// Используем внутренний db из Tx
		return tx
	}
	return s.db
}

func (s *Storage) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	// Начало транзакции
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Отложенный вызов Rollback в случае ошибки
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				s.logger.Errorf("failed to rollback transaction: %v", err)
			}
		}
	}()

	ctxTx := context.WithValue(ctx, txKey{}, tx)
	if tErr := txFunc(ctxTx); tErr != nil {
		return fmt.Errorf("run transaction: %w", tErr)
	}

	// Фиксация транзакции
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

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
