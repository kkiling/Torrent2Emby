package deliverystate

type ChoseTorrentOptions struct {
	// Пользователь выбрал конкретный торрента файл
	Href *string
	// Пользователь поменял поисковый запрос
	NewSearchQuery *string
}
