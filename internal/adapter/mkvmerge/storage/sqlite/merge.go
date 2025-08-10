package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkiling/goplatform/storagebase"

	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
)

func (s *Storage) Create(ctx context.Context, create *mkvmerge.CreateMergeResult) error {
	paramsJSON, err := json.Marshal(create.Params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	_, err = s.base.Next(ctx).ExecContext(ctx, `
        INSERT INTO mkv_merge (id, idempotency_key, params, status, created_at)
        VALUES (?, ?, ?, ?, ?)
    `, create.ID, create.IdempotencyKey, paramsJSON, create.Status, create.CreatedAt)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id uuid.UUID, update *mkvmerge.UpdateMergeResult) error {
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

	_, err := s.base.Next(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) UpdateProgress(ctx context.Context, id uuid.UUID, progress float64) error {
	query := "UPDATE mkv_merge SET progress = ? WHERE id = ?"
	_, err := s.base.Next(ctx).ExecContext(ctx, query, progress, id)
	if err != nil {
		return s.base.HandleError(err)
	}
	return nil
}

func (s *Storage) getMergeResult(row *sql.Row) (*mkvmerge.MergeResult, error) {
	var result mkvmerge.MergeResult
	var paramsJSON string
	var errorStr sql.NullString
	var completedAt sql.NullTime
	var progress sql.NullFloat64

	err := row.Scan(
		&result.ID,
		&paramsJSON,
		&result.Status,
		&errorStr,
		&result.CreatedAt,
		&completedAt,
		&progress,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storagebase.ErrNotFound
		}
		return nil, s.base.HandleError(err)
	}

	// Десериализуем параметры
	if err = json.Unmarshal([]byte(paramsJSON), &result.Params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal params: %w", err)
	}

	// Обрабатываем nullable поля
	if errorStr.Valid {
		result.Error = &errorStr.String
	}
	if completedAt.Valid {
		result.CompletedAt = &completedAt.Time
	}
	if progress.Valid {
		result.Progress = &progress.Float64
	}

	return &result, nil
}

func (s *Storage) GetByID(ctx context.Context, id uuid.UUID) (*mkvmerge.MergeResult, error) {
	row := s.base.Next(ctx).QueryRowContext(ctx, `
        SELECT id, params, status, error, created_at, completed_at, progress
        FROM mkv_merge WHERE id = ?
    `, id)
	return s.getMergeResult(row)
}

func (s *Storage) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*mkvmerge.MergeResult, error) {
	row := s.base.Next(ctx).QueryRowContext(ctx, `
        SELECT id, params, status, error, created_at, completed_at, progress
        FROM mkv_merge WHERE idempotency_key = ?
    `, idempotencyKey)
	return s.getMergeResult(row)
}

func (s *Storage) GetOldestUncompleted(ctx context.Context) (*mkvmerge.MergeResult, error) {
	row := s.base.Next(ctx).QueryRowContext(ctx, `
        SELECT id, params, status, error, created_at, completed_at, progress
        FROM mkv_merge
        WHERE completed_at is null
        ORDER BY created_at
        LIMIT 1
    `)
	return s.getMergeResult(row)
}

func (s *Storage) AddMergeLogs(ctx context.Context, id uuid.UUID, log mkvmerge.MergeLogs) error {
	_, err := s.base.Next(ctx).ExecContext(ctx, `
        INSERT INTO mkv_merge_logs (merge_id, type, content, created_at)
        VALUES (?, ?, ?, ?)
    `, id, log.Type, log.Content, log.CreatedAt)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) DeleteLogs(ctx context.Context, mergeID uuid.UUID) error {
	_, err := s.base.Next(ctx).ExecContext(ctx, `
        DELETE from mkv_merge_logs WHERE merge_id = ?
    `, mergeID)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}
