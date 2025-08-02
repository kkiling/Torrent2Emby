package delivery

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

type EpisodeInfo struct {
	// Номер сезона
	SeasonNumber int
	// Наименования эпизода
	EpisodeName string
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
	Episode    EpisodeInfo
	Video      VideoFile
	AudioFiles []Track
	Subtitles  []Track
}

// ---

type TorrentDownloadStatus struct {
	Progress   float64
	IsComplete bool
}

type CatalogsInfo struct {
	// Путь до каталога сериала
	TvShowPath string
	// Путь до каталога сезона
	TvShowSeasonPath string
}
