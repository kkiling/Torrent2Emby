package mkvmerge

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase"
	"github.com/samber/lo"
)

const retryDelay = time.Second * 5

type Pipeline struct {
	logger  log.Logger
	merger  MkvMerge
	storage Storage
}

func NewPipeline(merger MkvMerge, storage Storage, logger log.Logger) *Pipeline {
	return &Pipeline{
		logger:  logger.Named("mkv_merge_convener"),
		merger:  merger,
		storage: storage,
	}
}

func (s *Pipeline) AddToMerge(ctx context.Context, idempotencyKey string, params MergeParams) (*MergeResult, error) {
	if find, err := s.storage.GetByIdempotencyKey(ctx, idempotencyKey); err != nil {
		switch {
		case errors.Is(err, storagebase.ErrNotFound):
		default:
			return nil, fmt.Errorf("storage.GetByIdempotencyKey: %w", err)
		}
	} else if find != nil {
		return find, ErrAlreadyExists
	}

	result := MergeResult{
		ID:        uuid.New(),
		Params:    params,
		Status:    PendingStatus,
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

func (s *Pipeline) GetMergeResult(ctx context.Context, id uuid.UUID) (*MergeResult, error) {
	result, err := s.storage.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storagebase.ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("storage.GetFirstUncompletedMergeResult: %w", err)
		}
	}
	return result, nil
}

func (s *Pipeline) startTimer(ctx context.Context) error {
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

func (s *Pipeline) getProgress(content string) *float64 {
	// Регулярное выражение для поиска прогресса в формате "Progress: XX%"
	re := regexp.MustCompile(`Progress:\s*(\d+)%`)

	// Ищем совпадение в content
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil // не найдено совпадений
	}

	// Преобразуем найденное число в int
	progressInt, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil // не удалось преобразовать в число
	}

	// Проверяем, что прогресс в допустимом диапазоне (0-100)
	if progressInt < 0 {
		progressInt = 0
	} else if progressInt > 100 {
		progressInt = 100
	}

	// Конвертируем в float64 и делим на 100 для получения значения 0-1
	progress := float64(progressInt) / 100.0
	return &progress
}

func (s *Pipeline) runMerge(ctx context.Context, id uuid.UUID, params MergeParams) error {
	mergeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var outputChan = make(chan OutputMessage)
	// Закрываем канал при завершении функции
	defer close(outputChan)

	go func() {
		for msg := range outputChan {
			if msg.Type == ErrorMessageType {
				s.logger.Errorf("mvk merge logs: %s", msg.Content)
			} else {
				s.logger.Debugf("mvk merge logs: %s", msg.Content)
				progress := s.getProgress(msg.Content)
				if progress != nil {
					errUpdate := s.storage.UpdateProgress(ctx, id, *progress)
					if errUpdate != nil {
						s.logger.Errorf("mvk UpdateProgress: %s", errUpdate.Error())
					}
				}
			}

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

func (s *Pipeline) StartMergePipeline(ctx context.Context) error {
	for ctx.Err() == nil {
		result, err := s.storage.GetOldestUncompleted(ctx)
		if err != nil {
			switch {
			case errors.Is(err, storagebase.ErrNotFound):
				// Тут  таймер  - что бы не спамить базу
				if err = s.startTimer(ctx); err != nil {
					return err
				}
				continue
			default:
				return fmt.Errorf("storage.GetFirstUncompletedMergeResult: %w", err)
			}
		}

		s.logger.Debugf("start mvk merge: %s", result.Params.VideoInputFile)

		// TODO: транзакция
		err = s.storage.DeleteLogs(ctx, result.ID)
		if err != nil {
			return fmt.Errorf("storage.DeleteLogs: %w", err)
		}
		err = s.storage.Update(ctx, result.ID, &UpdateMergeResult{
			Status: RunningStatus,
		})
		if err != nil {
			return fmt.Errorf("storage.Update: %w", err)
		}

		err = s.runMerge(ctx, result.ID, result.Params)
		if err != nil {
			s.logger.Errorf("error mvk merge: %s", result.Params.VideoInputFile)

			uerr := s.storage.Update(ctx, result.ID, &UpdateMergeResult{
				Status:    ErrorStatus,
				Completed: lo.ToPtr(time.Now()),
				Error:     lo.ToPtr(err.Error()),
			})
			if uerr != nil {
				return fmt.Errorf("storage.Update: %w", uerr)
			}
			return fmt.Errorf("runMerge: %w", err)
		}

		s.logger.Debugf("complete mvk merge: %s", result.Params.VideoInputFile)

		err = s.storage.Update(ctx, result.ID, &UpdateMergeResult{
			Status:    CompleteStatus,
			Completed: lo.ToPtr(time.Now()),
		})

		if err != nil {
			return fmt.Errorf("storage.Update: %w", err)
		}

	}

	return ctx.Err()
}
