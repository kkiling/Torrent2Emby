package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/usercase/mkvmergepipeline"
	"github.com/kkiling/torrent2emby/internal/usercase/mkvmergepipeline/storage"
)

func (s *Storage) Create(ctx context.Context, create *mkvmergepipeline.CreateMergeResult) error {
	paramsJSON, err := json.Marshal(create.Params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	_, err = s.next(ctx).ExecContext(ctx, `
        INSERT INTO mkv_merge (id, idempotency_key, params, status, created_at)
        VALUES (?, ?, ?, ?, ?)
    `, create.ID, create.IdempotencyKey, paramsJSON, create.Status, create.CreatedAt)

	if err != nil {
		return handleError(err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id uuid.UUID, update *mkvmergepipeline.UpdateMergeResult) error {
	query := "UPDATE mkv_merge SET status = ?"
	args := []interface{}{update.Status}

	if update.Error != nil {
		query += ", error = ?"
		args = append(args, *update.Error)
	}

	if update.Completed != nil {
		query += ", completed_at = ?"
		args = append(args, *update.Completed)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := s.next(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func (s *Storage) getMergeResult(row *sql.Row) (*mkvmergepipeline.MergeResult, error) {
	var result mkvmergepipeline.MergeResult
	var paramsJSON string
	var errorStr sql.NullString
	var completedAt sql.NullTime

	err := row.Scan(
		&result.ID,
		&paramsJSON,
		&result.Status,
		&errorStr,
		&result.CreatedAt,
		&completedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, handleError(err)
	}

	// Десериализуем параметры
	if err := json.Unmarshal([]byte(paramsJSON), &result.Params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal params: %w", err)
	}

	// Обрабатываем nullable поля
	if errorStr.Valid {
		result.Error = &errorStr.String
	}
	if completedAt.Valid {
		result.CompletedAt = &completedAt.Time
	}

	return &result, nil
}

func (s *Storage) GetByID(ctx context.Context, id uuid.UUID) (*mkvmergepipeline.MergeResult, error) {
	row := s.next(ctx).QueryRowContext(ctx, `
        SELECT id, params, status, error, created_at, completed_at
        FROM mkv_merge WHERE id = ?
    `, id)
	return s.getMergeResult(row)
}

func (s *Storage) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*mkvmergepipeline.MergeResult, error) {
	row := s.next(ctx).QueryRowContext(ctx, `
        SELECT id, params, status, error, created_at, completed_at
        FROM mkv_merge WHERE idempotency_key = ?
    `, idempotencyKey)
	return s.getMergeResult(row)
}

func (s *Storage) GetOldestUncompleted(ctx context.Context) (*mkvmergepipeline.MergeResult, error) {
	row := s.next(ctx).QueryRowContext(ctx, `
        SELECT id, params, status, error, created_at, completed_at
        FROM mkv_merge
        WHERE completed_at is null
        ORDER BY created_at
        LIMIT 1
    `)

	return s.getMergeResult(row)
}

func (s *Storage) AddMergeLogs(ctx context.Context, id uuid.UUID, log mkvmergepipeline.MergeLogs) error {
	_, err := s.next(ctx).ExecContext(ctx, `
        INSERT INTO mkv_merge_logs (merge_id, type, content, created_at)
        VALUES (?, ?, ?, ?)
    `, id, log.Type, log.Content, log.CreatedAt)

	if err != nil {
		return handleError(err)
	}

	return nil
}

func (s *Storage) DeleteLogs(ctx context.Context, mergeID uuid.UUID) error {
	_, err := s.next(ctx).ExecContext(ctx, `
        DELETE from mkv_merge_logs WHERE merge_id = ?
    `, mergeID)

	if err != nil {
		return handleError(err)
	}

	return nil
}

/*
	query := "UPDATE mkv_merge SET status = ?"
	args := []interface{}{update.Status}

	if update.Error != nil {
		query += ", error = ?"
		args = append(args, *update.Error)
	}

	if update.Completed != nil {
		query += ", completed_at = ?"
		args = append(args, *update.Completed)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := s.next(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return handleError(err)
	}

	return nil
}
*/
