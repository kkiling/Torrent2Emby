package qbittorrent

// TorrentAddOptions содержит параметры для добавления нового торрента
type TorrentAddOptions struct {
	Magnet   string
	SavePath string
	Category string
	Tags     []string
	Paused   bool
}
