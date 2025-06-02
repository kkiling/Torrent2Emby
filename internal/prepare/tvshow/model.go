package tvshow

type Episode struct {
	// Номер эпизода в сезоне
	EpisodeNumber int
	// Наименование эпизода
	Name string
}

type TorrentFile struct {
	// Путь до файла торрента (путь относительно ContentPath)
	RelativePath string
	// Размер файла в байтах
	Size int64
}

type PrepareAudio struct {
	Name string
	File TorrentFile
}

type PrepareSubtitles struct {
	Name string
	File TorrentFile
}

type PrepareVideo struct {
	File TorrentFile
}

type PrepareEpisode struct {
	Episode    Episode
	VideoFile  *PrepareVideo
	AudioFiles []PrepareAudio
	Subtitles  []PrepareSubtitles
}

type PrepareTVShowSeason struct {
	Episodes []PrepareEpisode
}
