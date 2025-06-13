package videocontent

import (
	"context"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/samber/lo"
)

type Config struct {
	TVShowTorrentSavePath string
	BasePath              string // "/nfs"
}

type Service struct {
	config        Config
	repository    Repository
	tvShowLibrary TVShowLibrary
	torrentSite   TorrentSite
	torrentClient TorrentClient
	prepareTVShow PrepareTVShow
}

func NewService(
	config Config,
	repository Repository,
	tvShowLibrary TVShowLibrary,
	torrentSite TorrentSite,
	torrentClient TorrentClient,
	prepareTVShow PrepareTVShow,
) *Service {
	return &Service{
		repository:    repository,
		tvShowLibrary: tvShowLibrary,
		torrentSite:   torrentSite,
		torrentClient: torrentClient,
		prepareTVShow: prepareTVShow,
	}
}

// CreateNewVideoContent создание нового видео контента для сезона сериала / фильма из библиотеки
func (s *Service) CreateNewVideoContent(ctx context.Context, params CreateNewVideoContentParams) (*VideoContent, error) {
	if params.MediaID.MovieID == nil && params.MediaID.TVShow == nil {
		return nil, fmt.Errorf("movieID or TVShow is required: %w", ucerr.InvalidArgument)
	}
	if params.MediaID.MovieID != nil && params.MediaID.TVShow == nil {
		return nil, fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	// Проверка на то что у текущего VideoContent нет VideoContent
	if contents, err := s.repository.GetVideoContentByMediaID(ctx, params.MediaID); err != nil {
		return nil, fmt.Errorf("repository.GetVideoContentByMediaID: %w", err)
	} else if len(contents) > 0 {
		return nil, fmt.Errorf("video content for this media already exists: %w", ucerr.InvalidArgument)
	}

	content := VideoContent{
		ID:             0, // TODO: Генерация ID
		MediaID:        params.MediaID,
		DeliveryStatus: SearchTorrentsStatus, // Переходим на шаг поиска торрент раздачи
	}

	// Сохранение content в базу
	if err := s.repository.SaveVideoContent(ctx, &content); err != nil {
		return nil, fmt.Errorf("repository.SaveVideoContent: %w", err)
	}

	return &content, nil
}
func (s *Service) Complete(ctx context.Context, content *VideoContent) error {
	// TODO: крутим стейт машину
	panic("implement me")
}

func (s *Service) getTVShowQuery(ctx context.Context, tvShowID uint64, seasonNumber int) (string, error) {
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: tvShowID,
	})
	if err != nil {
		return "", fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return "", fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	season, find := lo.Find(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == seasonNumber
	})
	if !find {
		return "", fmt.Errorf("season not found: %w", ucerr.NotFound)
	}
	// Формируем поисковый запрос на основе инфы  о сезоне сериала
	searchQuery := fmt.Sprintf("%s сезон %d", tvShowInfo.Result.Name, season.SeasonNumber)

	return searchQuery, nil
}

// generateSearchQuery Автоматическое формирование поискового запроса
func (s *Service) generateSearchQuery(ctx context.Context, params GenerateSearchQueryParams) (string, error) {
	searchQuery := ""
	if params.MediaID.TVShow != nil {
		var err error
		searchQuery, err = s.getTVShowQuery(ctx, params.MediaID.TVShow.TVShowID, params.MediaID.TVShow.SeasonNumber)
		if err != nil {
			return "", fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
		}
	}
	if params.MediaID.MovieID != nil {
		return "", fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	return searchQuery, nil
}

// searchTorrent поиск раздачи
func (s *Service) searchTorrent(ctx context.Context, params SearchTorrentParams) (*TorrentSearchResult, error) {
	// Делаем запрос к торрент сайту, получаем список раздач
	searchResult, err := s.torrentSite.SearchTorrents(params.SearchQuery)
	if err != nil {
		return nil, fmt.Errorf("torrentSite.SearchTorrents: %w", err)
	}

	result := TorrentSearchResult{}
	for _, item := range searchResult.Results {
		result.Result = append(result.Result, TorrentSearch{
			Title:     item.Title,
			Href:      item.Href,
			Size:      item.Size,
			Seeds:     item.Seeds,
			Leeches:   item.Leeches,
			Downloads: item.Downloads,
			AddedDate: item.AddedDate,
		})
	}

	return &result, nil
}

// ChangeSearchQuery пользователь меняет поисковый запрос
func (s *Service) ChangeSearchQuery(ctx context.Context, params ChangeSearchQueryParams) error {
	// TODO: Обновление SearchQuery в базе
	// TODO: возвращаемся на шаг searchTorrent
	panic("implement me")
}

// ChoseTorrent пользователь выбирает конкретную раздачу для скачивания
func (s *Service) ChoseTorrent(ctx context.Context, params *ChoseTorrentParams) error {
	// TODO: сохраняем выбранный результат пользователя в базу
	// TODO: переходим на следующий шаг getMagnetLink
	panic("implement me")
}

// getMagnetLink получение магнет ссылки на основе выбора раздачи пользователем
func (s *Service) getMagnetLink(ctx context.Context, params GetMagnetLinkParams) (*MagnetInfo, error) {
	// Получение магнет ссылки
	magnetInfo, err := s.torrentSite.GetMagnetLink(params.Href)
	if err != nil {
		return nil, fmt.Errorf("torrentSite.GetMagnetLink: %w", err)
	}

	return &MagnetInfo{
		Magnet: magnetInfo.Magnet,
		Hash:   magnetInfo.Hash,
	}, nil
}

// CreateTorrent добавление торрент раздачи в торрент клиент
func (s *Service) createTorrent(ctx context.Context, params *CreateTorrentParams) error {
	if params.MediaID.MovieID != nil && params.MediaID.TVShow == nil {
		return fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	// Создание раздачи в торрент клиенте, выставление его сразу в паузу
	err := s.torrentClient.AddTorrent(qbittorrent.TorrentAddOptions{
		Magnet:   params.Magnet,
		SavePath: s.config.TVShowTorrentSavePath,
		Category: "tvshow",
		Tags: []string{
			fmt.Sprintf("tvshowID:%d", params.MediaID.TVShow.TVShowID),
			fmt.Sprintf("seasonNumber:%d", params.MediaID.TVShow.SeasonNumber),
		},
		Paused: false,
	})
	if err != nil {
		return fmt.Errorf("torrentClient.AddTorrent: %w", err)
	}

	return nil
}

// WaitingTorrentDownloadComplete ожидание завершения окончания скачивания раздачи
func (s *Service) waitingTorrentDownloadComplete(ctx context.Context, params WaitingTorrentDownloadCompleteParams) (*TorrentDownloadStatus, error) {
	// Достаем инфу о торрент раздаче
	torrentInfo, err := s.torrentClient.GetTorrentInfo(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if torrentInfo == nil {
		return nil, fmt.Errorf("torrentInfo not found: %w", ucerr.NotFound)
	}

	switch torrentInfo.State {
	case qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP:
		return &TorrentDownloadStatus{
			Progress:   torrentInfo.Progress,
			IsComplete: true,
		}, nil
	default:
		return &TorrentDownloadStatus{
			Progress:   torrentInfo.Progress,
			IsComplete: false,
		}, nil
	}
}

// MergeVideoFiles запуск обработки видеофайлов
func (s *Service) mergeVideoFiles(ctx context.Context, content *VideoContent) error {
	// TODO: формирование каталогов сериала на медиасервер
	// TODO: на основе FileMatches запуск mkvmerge сразу сохранением файлов в каталогах медиа сервера
	// TODO: Переход на следующий шаг - установки методаных
	panic("implement me")
}

// CopyFilesToMediaServer шаг копирования файлов на медиа сервер
func (s *Service) copyFilesToMediaServer(ctx context.Context, content *VideoContent) error {
	// TODO: формирование каталогов сериала на медиасервер
	// TODO: создание симлинков видеофайлов с торрент раздачи в каталогах медиасервера
	// TODO: Переход на следующий шаг - установки методаных
	panic("implement me")
}

// SetMetadataInMediaServer установка методанных
func (s *Service) setMetadataInMediaServer(ctx context.Context, content *VideoContent) error {
	// TODO: на основании информации о фильме/сереале
	// TODO: в emby устанавливаем методанные
	// TODO: Переход на следующий шаг
	panic("implement me")
}

// SendNotificationSuccessDelivery уведомление о успешной доставке
func (s *Service) sendNotificationSuccessDelivery(ctx context.Context, content *VideoContent) error {
	// TODO: уведомление о успешной доставке
	// TODO: Завершение доставки
	panic("implement me")
}
