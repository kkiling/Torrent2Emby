package rutracker

import "fmt"

var (
	// NotAuthorizedErr в момент запроса оказалось что клиент не авторизован, скорее всего протухли куки
	// Что бы исправить достаточно еще раз попробовать выполнить запрос
	NotAuthorizedErr = fmt.Errorf("not authorized")
	// AuthenticationFailedErr во время попытки логина на сайте произошла ошибка
	// Либо указан неверный пароль, либо сайт просит ввести капчу
	AuthenticationFailedErr = fmt.Errorf("authentication failed")
	// ServiceUnavailableErr сервис недоступен, можно попробовать повторить попытку
	ServiceUnavailableErr = fmt.Errorf("service unavailable")
)

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
