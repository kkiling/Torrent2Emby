package content

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase"
	"github.com/kkiling/statemachine"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners/tvshowdeliverystate"
	"github.com/samber/lo"

	"time"
)

type Service struct {
	logger              log.Logger
	storage             Storage
	tvShowLibrary       TVShowLibrary
	tvShowDeliveryState TVShowDeliveryState
	clock               Clock
	uuidGenerator       UUIDGenerator
}

func NewService(
	logger log.Logger,
	storage Storage,
	tvShowLibrary TVShowLibrary,
	tvShowDeliveryState TVShowDeliveryState,
) *Service {
	return &Service{
		storage:             storage,
		tvShowLibrary:       tvShowLibrary,
		tvShowDeliveryState: tvShowDeliveryState,
		clock:               &realClock{},
		uuidGenerator:       &uuidGenerator{},
		logger:              logger.Named("content"),
	}
}

/*
	Какие кейсы:
		- Создание новой доставки файлов с прохождением полного флоу от поиска раздачи до доставки файлов до медиасервера
		- Можно оставить информацию о раздаче (Href и Magnet) но при этом удалить все файлы, что бы не занимали место на диске
		- Потом на основе (Href и Magnet) восстанавливать файлы и скачивать их снова, при этом не запрашивая больше инфу от клиента
			и все подтягивать из старых стейтов (что делать если раздача обновиться?)
       - Раздача может обновиться и запускается процесс обновления раздачи
*/

// CreateVideoContent создание файловой раздачи
func (s *Service) CreateVideoContent(ctx context.Context, params CreateVideoContentParams) (*VideoContent, error) {
	found, err := s.GetVideoContent(ctx, params.ContentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}
	if len(found) > 0 {
		return nil, ucerr.AlreadyExists
	}

	if params.ContentID.MovieID != nil {
		return nil, fmt.Errorf("movieID is not support: %w", ucerr.InvalidArgument)
	}

	// Добавить проверку на наличие уже раздачи и выкидывать ее
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: params.ContentID.TVShow.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return nil, fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	// TODO: !!! !!! !!! подумать как обернуть в одну транзакцию
	state, err := s.tvShowDeliveryState.Create(ctx, tvshowdeliverystate.CreateOptions{
		TVShowID: *params.ContentID.TVShow,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.Create: %w", err)
	}

	videoContent := VideoContent{
		ID:             s.uuidGenerator.New(),
		CreatedAt:      s.clock.Now(),
		ContentID:      params.ContentID,
		DeliveryStatus: DeliveryStatusInProgress,
		State: []State{
			{
				StateID: state.ID,
				Type:    runners.TVShowDelivery,
			},
		},
	}

	if err = s.storage.SaveVideoContent(ctx, &videoContent); err != nil {
		return nil, fmt.Errorf("storage.SaveVideoContent: %w", err)
	}

	return &videoContent, nil
}

func (s *Service) GetVideoContent(ctx context.Context, contentID videocontent.ContentID) ([]VideoContent, error) {
	result, err := s.storage.GetVideoContent(ctx, contentID)
	switch {
	case err == nil:
	case errors.Is(err, storagebase.ErrNotFound): // Выпуск не найден
		return nil, ucerr.NotFound
	default:
		return nil, fmt.Errorf("storage.GetVideoContent: %w", err)
	}
	return result, nil
}

func (s *Service) getStateID(ctx context.Context, contentID videocontent.ContentID, runersType runners.Type) (uuid.UUID, error) {
	contents, err := s.storage.GetVideoContent(ctx, contentID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("storage.GetVideoContent: %w", err)
	}

	if len(contents) != 1 {
		return uuid.UUID{}, ucerr.NotFound
	}

	content := contents[0]
	state, find := lo.Find(content.State, func(item State) bool {
		return item.Type == runersType
	})
	if !find {
		return uuid.UUID{}, ucerr.NotFound
	}

	return state.StateID, nil
}

func (s *Service) GetTVShowDeliveryData(ctx context.Context, contentID videocontent.ContentID) (*tvshowdeliverystate.TVShowDeliveryData, error) {
	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelivery)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}

	result, err := s.tvShowDeliveryState.GetStateByID(ctx, stateID)
	if err != nil {
		return nil, fmt.Errorf("s.GetStateByID: %w", err)
	}

	return &result.Data, nil
}

func (s *Service) ChoseTorrentOptions(ctx context.Context,
	contentID videocontent.ContentID,
	opts tvshowdeliverystate.ChoseTorrentOptions,
) (*tvshowdeliverystate.TVShowDeliveryData, error) {
	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelivery)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}
	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateID, opts)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
	}
	if executeErr != nil {
		return nil, executeErr
	}
	return &newState.Data, nil
}

func (s *Service) ChoseFileMatchesOptions(ctx context.Context,
	contentID videocontent.ContentID,
	opts tvshowdeliverystate.ChoseFileMatchesOptions,
) (*tvshowdeliverystate.TVShowDeliveryData, error) {
	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelivery)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}
	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateID, opts)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
	}
	if executeErr != nil {
		return nil, executeErr
	}
	return &newState.Data, nil
}

func (s *Service) completeTVShowDelivery(ctx context.Context, content VideoContent) error {
	stateInfo, find := lo.Find(content.State, func(item State) bool {
		return item.Type == runners.TVShowDelivery
	})
	if !find {
		return ucerr.NotFound
	}
	s.logger.Infof("start complete content: %d (season %d) - status %s",
		content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber, content.DeliveryStatus)

	if content.DeliveryStatus == DeliveryStatusInProgress {
		// Добиваем стейт
		newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateInfo.StateID)
		if err != nil && !errors.Is(err, statemachine.ErrOptionsIsUndefined) {
			return fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
		}
		if executeErr != nil {
			s.logger.Infof("executeError: %v", executeErr)
			return executeErr
		}

		s.logger.Infof("newState step: %s", newState.Step)

		needUpdate := false
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatus,
		}
		if newState.Status == statemachine.CompletedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusDelivered
		} else if newState.Status != statemachine.FailedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusFailed
		} else if content.TorrentInfo == nil && newState.Data.SelectTorrentHref != nil && newState.Data.MagnetInfo != nil {
			needUpdate = true
			updateVideoContent.TorrentInfo = &TorrentInfo{
				Href:   *newState.Data.SelectTorrentHref,
				Magnet: newState.Data.MagnetInfo.Magnet,
				Hash:   newState.Data.MagnetInfo.Hash,
			}
		}

		if needUpdate {
			s.logger.Infof("update content: %d (season %d)",
				content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber)
			err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent)
			if err != nil {
				return fmt.Errorf("storage.UpdateVideoContent: %w", err)
			}
		}

	}

	//Трекаем обновление статуса в процессе доставки in_progress до доставлено delivered на основе стейта
	//	in_progress -> delivered
	//Трекаем по аналогии
	//	updating -> delivered
	//Трекаем по анлогии на основании стейта
	//	deleting -> deleted

	return nil
}

func (s *Service) completeTVShowDeliveries(ctx context.Context) error {
	contents, err := s.storage.GetVideoContents(ctx, DeliveryStatusInProgress, 10)
	if err != nil {
		return fmt.Errorf("storage.GetVideoContents: %w", err)
	}
	for _, content := range contents {
		err = s.completeTVShowDelivery(ctx, content)
		if err != nil {
			return fmt.Errorf("completeTVShowDelivery: %w", err)
		}
	}
	return nil
}

func (s *Service) Complete(ctx context.Context) error {
	scheduler := gocron.NewScheduler(time.UTC)

	// Настраиваем выполнение в 1 поток (по умолчанию и так последовательно)
	scheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)

	// Запускаем задачу каждые 3 секунды
	_, err := scheduler.Every(3).Seconds().Do(func() {
		select {
		case <-ctx.Done(): // Если контекст отменён, выходим
			return
		default:
			if err := s.completeTVShowDeliveries(ctx); err != nil {
				s.logger.Errorf("completeTVShowDeliveries: %v", err)
			}
		}
	})
	if err != nil {
		return fmt.Errorf("failed run cron complete: %w", err)
	}

	// Запускаем планировщик (асинхронно)
	scheduler.StartAsync()

	// Ждём отмены контекста
	<-ctx.Done()

	// Останавливаем планировщик при завершении
	scheduler.Stop()
	return nil
}
