package tvshow

type Episode struct {
	// Номер эпизода в сезоне
	EpisodeNumber int
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
	EpisodeNumber int
	VideoFile     *PrepareVideo
	AudioFiles    []PrepareTrack
	Subtitles     []PrepareTrack
}

type PrepareTVShowSeason struct {
	Episodes []PrepareEpisode
}
