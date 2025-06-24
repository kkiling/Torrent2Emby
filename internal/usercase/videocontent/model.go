package videocontent

type MediaID struct {
	MovieID *uint64
	TVShow  *struct {
		TVShowID     uint64
		SeasonNumber int
	}
}

// VideoContent модель содержащая информацию о видео контенте для сезона сериала / фильма
/*
	К одному сезону сериала / фильму может быть привязано несколько VideoContent,
	но для упрощения пока будем пока разрешать только 1
*/
type VideoContent struct {
	ID      uint64
	MediaID MediaID

	// Хеш торрент раздачи
	// Магнет ссылку раздачи
	// Ссылку на раздачу (ссылка с торрент трекера)

	// Каталог раздачи
	// Каталог в медиасервере
	// Статус доставки контента
	DeliveryStatus StatusVideoDelivery

	// Данные выпуска
	SearchQuery           *string
	TorrentSearch         *TorrentSearchResult
	SelectTorrentHref     *string
	MagnetInfo            *MagnetInfo
	ContentMatches        []ContentMatches
	TorrentDownloadStatus *TorrentDownloadStatus
	CatalogsInfo          *CatalogsInfo
	MergeVideoStatus      *MergeVideoStatus
}

// StatusVideoDelivery статус доставки видео файлов до медиа сервера
type StatusVideoDelivery string

const (
	// GenerateSearchQuery - генерация запросса к трекеру
	GenerateSearchQuery StatusVideoDelivery = "generate_search_query"
	// SearchTorrents - ищем раздачи сезона сериала / фильма
	SearchTorrents StatusVideoDelivery = "search_torrents"
	// WaitingUserChoseTorrent - ожидание когда пользователь выберет раздачу
	WaitingUserChoseTorrent StatusVideoDelivery = "waiting_user_chose_torrent"
	// GetMagnetLink получение магнет ссылки
	GetMagnetLink StatusVideoDelivery = "get_magnet_link_status"
	// AddTorrentToTorrentClient Добавление раздачи для скачивания торрент клиентом
	AddTorrentToTorrentClient StatusVideoDelivery = "add_torrent_client_status"
	// PrepareFileMatches получение информации о файлах раздачи
	PrepareFileMatches StatusVideoDelivery = "prepare_file_matches"
	// WaitingChoseFileMatches ожидание подтверждения пользователем соответствий выбора файлов
	WaitingChoseFileMatches StatusVideoDelivery = "waiting_chose_file_matches"
	// WaitingTorrentDownloadComplete ожидание завершения окончания скачивания раздачи
	WaitingTorrentDownloadComplete StatusVideoDelivery = "waiting_torrent_download_complete"

	// CreateVideoContentCatalogs Формирование каталогов и иерархии файлов
	CreateVideoContentCatalogs StatusVideoDelivery = "create_video_content_catalogs"
	// DeterminingNeedConvertFiles Определение необходимости конвертации файлов
	DeterminingNeedConvertFiles StatusVideoDelivery = "determining_need_convert_files"
	// --- Ветвь если необходимо добавление аудио дорожек/субтитров

	// MergeVideoFiles Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
	MergeVideoFiles StatusVideoDelivery = "merge_video_files"

	// -- Ветвь если не нужно изменять исходные файлы

	// CopyVideoFiles Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
	CopyVideoFiles StatusVideoDelivery = "copy_video_files"

	// SetMediaMetaData установка методаных серий сезона сериала / фильма в медиасервере
	SetMediaMetaData StatusVideoDelivery = "set_media_meta_data"

	// SendDeliveryNotification Отправка уведомления в telegramm о успешной доставки видеофайлов до медиа сервера
	SendDeliveryNotification StatusVideoDelivery = "send_delivery_notification"
)

type TorrentSearch struct {
	Title     string
	Href      string
	Size      string // TODO: переделать на int
	Seeds     string // TODO: переделать на int
	Leeches   string // TODO: переделать на int
	Downloads string // TODO: переделать на int
	AddedDate string // TODO: переделать на time.Time
}

type TorrentSearchResult struct {
	Result []TorrentSearch
}

type MagnetInfo struct {
	Magnet string
	Hash   string
}

// --- --- --- --- ---

type ContentInfo struct {
	// Наименования
	Name string
	// Постер и тд
}

type FileInfo struct {
	// Относительный путь до файла торрента (относительно каталога скачивания)
	RelativePath string
	// Полный путь до файла в системе
	FullPath string
	// Размер файла в байтах
	Size int64
	// Расширение файла
	Extension string
}

type VideoFile struct {
	File FileInfo
}

type Track struct {
	Name     string
	Language string
	File     FileInfo
}

// ContentMatches сопоставление видео файла с торрент файлом
type ContentMatches struct {
	ContentInfo ContentInfo
	Video       VideoFile
	AudioFiles  []Track
	Subtitles   []Track
}

// ---

type TorrentDownloadStatus struct {
	Progress   float64
	IsComplete bool
}

type MergeVideoStatus struct {
	// Сколько обработано файлов
	ProcessedFiles int
	// Сколько файлов всего нужно обработать
	AllFiles int
	// Конвертация завершена
	IsComplete bool
}

type CatalogsInfo struct {
	CatalogPath string
}
