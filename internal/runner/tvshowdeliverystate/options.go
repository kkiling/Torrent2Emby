package tvshowdeliverystate

type ChoseTorrentOptions struct {
	// Пользователь выбрал конкретный торрента файл
	Href *string
	// Пользователь поменял поисковый запрос
	NewSearchQuery *string
}

type ChoseFileMatchesOptions struct {
	// Пользователь подтверждает сметченные файлы
	Approve bool
}
