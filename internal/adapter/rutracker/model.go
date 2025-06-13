package rutracker

type Torrent struct {
	// Title наименование
	Title string
	// Href ссылка на раздачу
	Href string
	// Forum раздел раздачи
	Forum string
	// Author автор раздачи
	Author string
	// Size Размер раздачи
	Size string
	// Seeds Информация о сидах
	Seeds string
	// Leeches Информация о личах
	Leeches string
	// Downloads количество скачиваний
	Downloads string
	// AddedDate Дата добавления
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
