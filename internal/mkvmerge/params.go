package mkvmerge

type Track struct {
	Path     string
	Language string
	Name     string
	Default  bool
}

type MergeParams struct {
	VideoInputFile  string
	VideoOutputFile string
	AudioTracks     []Track
	SubtitleTracks  []Track
}
