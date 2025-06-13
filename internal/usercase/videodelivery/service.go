package videodelivery

import "context"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// StateVideoDelivery статус доставки видео файлов до медиа сервера
type StateVideoDelivery string

const (
	// NewState - раздача создана
	NewState StateVideoDelivery = "new"
	// SearchTorrentsState - ищем раздачи сезона сериала / фильма
	SearchTorrentsState StateVideoDelivery = "search_torrents"
	// WaitingUserChoseTorrentState - ожидание когда пользователь выберет раздачу
	WaitingUserChoseTorrentState StateVideoDelivery = "waiting_user_chose_torrent"
	// Добавление раздачи для скачивания торрент клиентом
	// Ожидание получения файлов раздачи
	// Формирование соответствия файлов раздачи с сезоном сериала / фильма
	// Ожидание когда клиент подтвердит/отредактирует соответствие файлов раздачи с сезоном сериала / фильма

	// --- Ветвь если необходимо добавление аудио дорожек/субтитров
	// Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера

	// -- Ветвь если не нужно изменять исходные файлы
	// Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)

	// Правка методаных серий сезона сериала / фильма в медиасервере
	// Отправка уведомления в telegramm о успешной доставки видеофайлов до медиа сервера

	// Доставлено
)

type VideoDelivery struct {
}

type CreateParams struct {
}

// Create создание доставки видео контента
func (s *Service) Create(ctx context.Context, createParams CreateParams) (*VideoDelivery, error) {
	// Создание доставки видео
	panic("implement me")
}

// Complete Стейт машина
func (s *Service) Complete(ctx context.Context, id uint64) error {
	//
	panic("implement me")
}
