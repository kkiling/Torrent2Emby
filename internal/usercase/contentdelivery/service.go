package contentdelivery

type Config struct {
	// BasePath Базовый путь от которого расположены все файлы торрента или медиа сервера
	// Например скачанные сериалы лежат по пути BasePath + TVShowTorrentSavePath
	BasePath string // "/nfs"
	// TVShowTorrentSavePath путь сохранения сериалов относительно торрент клиента
	TVShowTorrentSavePath string
	// TvShowMediaSavePath путь сохранения сериалов относительно медиа сервера
	TvShowMediaSavePath string
}

type Service struct {
	config        Config
	tvShowLibrary TVShowLibrary
	torrentSite   TorrentSite
	torrentClient TorrentClient
	prepareTVShow PrepareTVShow
	mkvMerge      MkvMerge
}

func NewService(
	config Config,
	tvShowLibrary TVShowLibrary,
	torrentSite TorrentSite,
	torrentClient TorrentClient,
	prepareTVShow PrepareTVShow,
	mkvMerge MkvMerge,
) *Service {
	return &Service{
		config:        config,
		tvShowLibrary: tvShowLibrary,
		torrentSite:   torrentSite,
		torrentClient: torrentClient,
		prepareTVShow: prepareTVShow,
		mkvMerge:      mkvMerge,
	}
}

/*
// Complete управление стейт машиной
func (s *Service) Complete(ctx context.Context, content *VideoContent) error {
	switch content.DeliveryStatus {
	case GenerateSearchQuery:
		// Генерация запроса
		result, err := s.generateSearchQuery(ctx, GenerateSearchQueryParams{
			MediaID: content.MediaID,
		})
		if err != nil {
			return fmt.Errorf("generateSearchQuery: %w", err)
		}
		content.SearchQuery = &result
		content.DeliveryStatus = SearchTorrents

	case SearchTorrents:
		// ищем раздачи сезона сериала / фильма
		result, err := s.searchTorrent(ctx, SearchTorrentParams{
			SearchQuery: *content.SearchQuery,
		})
		if err != nil {
			return fmt.Errorf("searchTorrent: %w", err)
		}
		content.TorrentSearch = result
		content.DeliveryStatus = WaitingUserChoseTorrent
	case WaitingUserChoseTorrent:
		// Ожидаем когда пользователь выберет раздачу
		// Или ожидаем что клиент изменит поисковый запрос, тогда прыгаем на SearchTorrents
		// TODO:
		content.SelectTorrentHref = lo.ToPtr("TODO get href")
		content.DeliveryStatus = GetMagnetLink
	case GetMagnetLink:
		// Получение магнет ссылки
		result, err := s.getMagnetLink(ctx, GetMagnetLinkParams{
			Href: *content.SelectTorrentHref,
		})
		if err != nil {
			return fmt.Errorf("searchTorrent: %w", err)
		}
		content.MagnetInfo = result
		content.DeliveryStatus = AddTorrentToTorrentClient
	case AddTorrentToTorrentClient:
		//  Добавление раздачи для скачивания торрент клиентом
		err := s.addTorrentToTorrentClient(ctx, AddTorrentParams{
			MediaID: content.MediaID,
			Magnet:  content.MagnetInfo.Magnet,
		})
		if err != nil {
			return fmt.Errorf("addTorrentToTorrentClient: %w", err)
		}
		content.DeliveryStatus = PrepareFileMatches
	case PrepareFileMatches:
		// Получение информации о файлах раздачи
		result, err := s.prepareFileMatches(ctx, PreparingFileMatchesParams{
			Hash:    content.MagnetInfo.Hash,
			MediaID: content.MediaID,
		})
		if err != nil {
			return fmt.Errorf("prepareFileMatches: %w", err)
		}
		if len(result) > 0 {
			content.ContentMatches = result
			content.DeliveryStatus = WaitingChoseFileMatches
		}
	case WaitingChoseFileMatches:
		// ожидание подтверждения пользователем соответствий выбора файлов
		// TODO:
		content.DeliveryStatus = WaitingTorrentDownloadComplete
	case WaitingTorrentDownloadComplete:
		// Ожидание когда торрент докачается до конца
		result, err := s.waitingTorrentDownloadComplete(ctx, WaitingTorrentDownloadCompleteParams{
			Hash: content.MagnetInfo.Hash,
		})
		if err != nil {
			return fmt.Errorf("waitingTorrentDownloadComplete: %w", err)
		}
		content.TorrentDownloadStatus = result
		if result.IsComplete {
			// Переход на следующий шаг
			content.DeliveryStatus = CreateVideoContentCatalogs
		}
	case CreateVideoContentCatalogs:
		// Формирование каталогов и иерархии файлов
		result, err := s.createVideoContentCatalogs(ctx, CreateVideoContentCatalogsParams{
			MediaID: content.MediaID,
		})
		if err != nil {
			return fmt.Errorf("waitingTorrentDownloadComplete: %w", err)
		}
		content.CatalogsInfo = &result
		content.DeliveryStatus = DeterminingNeedConvertFiles
	case DeterminingNeedConvertFiles:
		// Определение необходимости конвертации файлов
		needToMerge := false
		for _, m := range content.ContentMatches {
			// Если есть субтитры или аудиодорожки то нужно мержить
			if len(m.AudioFiles) > 0 || len(m.Subtitles) > 0 {
				needToMerge = true
				break
			}
		}
		if needToMerge {
			content.DeliveryStatus = MergeVideoFiles
		} else {
			content.DeliveryStatus = CopyVideoFiles
		}
	case CopyVideoFiles:
		// Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
		// TODO: CopyFiles
	case MergeVideoFiles:
		//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
		result, err := s.mergeVideoFiles(ctx, MergeVideoFilesParams{
			Hash:           content.MagnetInfo.Hash,
			ContentPath:    content.CatalogsInfo.CatalogPath,
			ContentMatches: content.ContentMatches,
			ProcessedFiles: func() int { // Стартуем с последнего
				if content.MergeVideoStatus != nil {
					return content.MergeVideoStatus.ProcessedFiles
				}
				return 0
			}(),
		})
		if err != nil {
			return fmt.Errorf("mergeVideoFiles: %w", err)
		}

		content.MergeVideoStatus = &result
		if result.IsComplete {
			// Переход на следующий шаг
			content.DeliveryStatus = SetMediaMetaData
		}
	case SetMediaMetaData:
		// установка методаных серий сезона сериала / фильма в медиасервере
		// TODO:
		content.DeliveryStatus = SendDeliveryNotification
	case SendDeliveryNotification:
		// TODO
		// Completed
	}

	return nil
}

// -- Реакция пользователя

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

// ChoseFileMatches подтверждение пользователем соответствия выбора файлов
func (s *Service) ChoseFileMatches(ctx context.Context, params ChoseFileMatchesParams) error {
	// TODO: Сохранение информации о соответствии файлов
	// TODO: Переход на следующий шаг
	panic("implement me")
}
*/
