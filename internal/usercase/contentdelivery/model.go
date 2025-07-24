package contentdelivery

type TVShowID struct {
	TVShowID     uint64
	SeasonNumber int
}
type MediaID struct {
	MovieID *uint64
	TVShow  *TVShowID
}

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
	// Номер сезона
	SeasonNumber int
	// Номер эпизода
	EpisodeNumber int
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
	TvShowCatalogPath string
}
