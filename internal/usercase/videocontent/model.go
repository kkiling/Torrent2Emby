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
}

// StatusVideoDelivery статус доставки видео файлов до медиа сервера
type StatusVideoDelivery string

const (
	// SearchTorrentsStatus - ищем раздачи сезона сериала / фильма
	SearchTorrentsStatus StatusVideoDelivery = "search_torrents"
	// WaitingUserChoseTorrentStatus - ожидание когда пользователь выберет раздачу
	WaitingUserChoseTorrentStatus StatusVideoDelivery = "waiting_user_chose_torrent"
	// Добавление раздачи для скачивания торрент клиентом
	// Ожидание получения файлов раздачи
	// Формирование соответствия файлов раздачи с сезоном сериала / фильма
	// Ожидание когда клиент подтвердит/отредактирует соответствие файлов раздачи с сезоном сериала / фильма

	// --- Ветвь если необходимо добавление аудио дорожек/субтитров
	// Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера

	// -- Ветвь если не нужно изменять исходные файлы
	// Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)

	// Правка методаных серий сезона сериала / фильма в медиасервере
	// Отправка уведомления в telegramm о успешной доставки видеофайлов до медиа сервера

	// Доставлено
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
