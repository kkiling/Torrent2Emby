package rutracker

type Torrent struct {
	Title     string
	Href      string
	Forum     string
	Author    string
	Size      string
	Seeds     string
	Leeches   string
	Downloads string
	AddedDate string
}

type TorrentResponse struct {
	Page         int
	TotalResults int
	Results      []Torrent
}

type MagnetInfo struct {
	Magnet string
	Hash   string
}
