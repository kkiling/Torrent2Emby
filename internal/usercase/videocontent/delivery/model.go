package delivery

import (
	"path/filepath"

	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
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

type TorrentInfo struct {
	Href   string
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

type TorrentState string

const (
	TorrentStateError       TorrentState = "error"       // Ошибка
	TorrentStateUploading   TorrentState = "uploading"   // Раздача (сидирование)
	TorrentStateDownloading TorrentState = "downloading" // Загрузка
	TorrentStateStopped     TorrentState = "stopped"     // Приостановлен (все виды паузы)
	TorrentStateQueued      TorrentState = "queued"      // В очереди
	TorrentStateUnknown     TorrentState = "unknown"     // Неизвестный статус
)

func mapTorrentState(qbState qbittorrent.TorrentState) TorrentState {
	switch qbState {
	// Ошибки
	case qbittorrent.TorrentStateError,
		qbittorrent.TorrentStateMissingFiles:
		return TorrentStateError

	// Раздача (сидирование)
	case qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP,
		qbittorrent.TorrentStateStalledUP,
		qbittorrent.TorrentStateCheckingUP,
		qbittorrent.TorrentStateQueuedUP,
		qbittorrent.TorrentStateForcedUP:
		return TorrentStateUploading

	// Загрузка
	case qbittorrent.TorrentStateDownloading,
		qbittorrent.TorrentStateMetaDL,
		qbittorrent.TorrentStateAllocating,
		qbittorrent.TorrentStateForcedDL,
		qbittorrent.TorrentStateMoving:
		return TorrentStateDownloading

	// Приостановлен
	case qbittorrent.TorrentStatePausedDL,
		qbittorrent.TorrentStateStoppedDL:
		return TorrentStateStopped

	// В очереди
	case qbittorrent.TorrentStateQueuedDL,
		qbittorrent.TorrentStateCheckingDL,
		qbittorrent.TorrentStateStalledDL,
		qbittorrent.TorrentStateCheckingResumeData:
		return TorrentStateQueued

	default:
		return TorrentStateUnknown
	}
}

type TorrentDownloadStatus struct {
	ContentPath string
	State       TorrentState
	Progress    float64
	IsComplete  bool
}

type TVShowCatalogPath struct {
	// Путь до каталога сериала
	TVShowPath string
	// Путь до каталога сезона (относительно каталога сериала)
	SeasonPath string
}

func (m TVShowCatalogPath) FullSeasonPath() string {
	return filepath.Join(m.TVShowPath, m.SeasonPath)
}

type TVShowCatalog struct {
	// Путь до раздачи сезона сериала
	TorrentPath string
	// Размер файлов раздачи сезона сериала (байты)
	TorrentSize uint64
	// Путь до сезона сериала на медиасервере
	MediaServerPath TVShowCatalogPath
	// Размер файлов раздачи сезона сериала (байты)
	MediaServerSize uint64
	// Файлы скопированы с раздачи или созданы ссылочная связь
	// True - файлы скопированы
	// False - файлы созданы через линки
	IsCopyFilesInMediaServer bool
}
