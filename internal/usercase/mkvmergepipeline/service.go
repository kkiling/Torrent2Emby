package mkvmergepipeline

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/usercase/mkvmergepipeline/storage"
	"github.com/samber/lo"
	"time"
)

const retryDelay = time.Second * 5

type Config struct {
}

type Service struct {
	cfg     Config
	logger  log.Logger
	merger  MkvMerge
	storage Storage
}

func NewService(cfg Config, merger MkvMerge, storage Storage, logger log.Logger) *Service {
	return &Service{
		cfg:     cfg,
		logger:  logger.Named("mkv_merge_convener"),
		merger:  merger,
		storage: storage,
	}
}

func (s *Service) AddToMerge(ctx context.Context, idempotencyKey string, params mkvmerge.MergeParams) (*MergeResult, error) {
	if find, err := s.storage.GetByIdempotencyKey(ctx, idempotencyKey); err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
		default:
			return nil, fmt.Errorf("storage.GetByIdempotencyKey: %w", err)
		}
	} else if find != nil {
		return find, ErrAlreadyExists
	}

	result := MergeResult{
		ID:        uuid.New(),
		Params:    params,
		Status:    Pending,
		CreatedAt: time.Now(),
	}

	if err := s.storage.Create(ctx, &CreateMergeResult{
		ID:             result.ID,
		IdempotencyKey: idempotencyKey,
		Params:         result.Params,
		Status:         result.Status,
		CreatedAt:      result.CreatedAt,
	}); err != nil {
		return nil, fmt.Errorf("storage.Create: %w", err)
	}

	return &result, nil
}

func (s *Service) GetMergeResult(ctx context.Context, id uuid.UUID) (*MergeResult, error) {
	result, err := s.storage.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("storage.GetFirstUncompletedMergeResult: %w", err)
		}
	}
	return result, nil
}

func (s *Service) startTimer(ctx context.Context) error {
	// Используем таймер с select для корректной обработки отмены контекста
	timer := time.NewTimer(retryDelay)
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (s *Service) runMerge(ctx context.Context, id uuid.UUID, params mkvmerge.MergeParams) error {
	mergeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var outputChan = make(chan mkvmerge.OutputMessage)
	// Закрываем канал при завершении функции
	defer close(outputChan)

	go func() {
		for msg := range outputChan {
			logErr := s.storage.AddMergeLogs(ctx, id, MergeLogs{
				CreatedAt: time.Now(),
				Type:      msg.Type,
				Content:   msg.Content,
			})
			if logErr != nil {
				s.logger.Errorf("AddMergeLogs: %v", logErr)
			}
		}
	}()

	if err := s.merger.Merge(mergeCtx, params, outputChan); err != nil {
		return fmt.Errorf("merger.Merge: %w", err)
	}
	return nil
}

func (s *Service) StartMergePipeline(ctx context.Context) error {
	for ctx.Err() == nil {
		result, err := s.storage.GetOldestUncompleted(ctx)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrNotFound):
				// Тут  таймер  - что бы не спамить базу
				if err = s.startTimer(ctx); err != nil {
					return err
				}
				continue
			default:
				return fmt.Errorf("storage.GetFirstUncompletedMergeResult: %w", err)
			}
		}

		// TODO: транзакция
		err = s.storage.DeleteLogs(ctx, result.ID)
		if err != nil {
			return fmt.Errorf("storage.DeleteLogs: %w", err)
		}
		err = s.storage.Update(ctx, result.ID, &UpdateMergeResult{
			Status: Running,
		})
		if err != nil {
			return fmt.Errorf("storage.Update: %w", err)
		}

		err = s.runMerge(ctx, result.ID, result.Params)
		if err != nil {
			uerr := s.storage.Update(ctx, result.ID, &UpdateMergeResult{
				Status:    Error,
				Completed: lo.ToPtr(time.Now()),
				Error:     lo.ToPtr(err.Error()),
			})
			if uerr != nil {
				return fmt.Errorf("storage.Update: %w", uerr)
			}
			return fmt.Errorf("runMerge: %w", err)
		}

		// TODO: научится определять ошибки обработки
		err = s.storage.Update(ctx, result.ID, &UpdateMergeResult{
			Status:    Complete,
			Completed: lo.ToPtr(time.Now()),
		})

		if err != nil {
			return fmt.Errorf("storage.Update: %w", err)
		}
	}

	return ctx.Err()
}
