package matchtvshow

type EpisodeInfo struct {
	// Номер сезона
	SeasonNumber int
	// Номер эпизода в сезоне
	EpisodeNumber int
	// Название эпизода, пока не нужно, но может будет нужно для метча в будущем
	EpisodeName string
}

type TorrentFile struct {
	// Путь до файла торрента (путь относительно ContentPath)
	RelativePath string
	//
	FullPath string
	// Размер файла в байтах
	Size int64
	// Расширение файла
	Extension string
}

type PrepareTrack struct {
	Name     string
	Language string
	File     TorrentFile
}

type PrepareVideo struct {
	File TorrentFile
}

type PrepareEpisode struct {
	Episode    EpisodeInfo
	VideoFile  *PrepareVideo
	AudioFiles []PrepareTrack
	Subtitles  []PrepareTrack
}

type PrepareTVShowSeason struct {
	Episodes []PrepareEpisode
}
