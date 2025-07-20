package deliverystate

import (
	"fmt"

	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
)

type Type string

const (
	ContentDeliveryType Type = "content_delivery"
)

// StepDelivery статус доставки видео файлов до медиа сервера
type StepDelivery string

const (
	// GenerateSearchQuery - генерация запросса к трекеру
	GenerateSearchQuery StepDelivery = "generate_search_query"
	// SearchTorrents - ищем раздачи сезона сериала / фильма
	SearchTorrents StepDelivery = "search_torrents"
	// WaitingUserChoseTorrent - ожидание когда пользователь выберет раздачу
	WaitingUserChoseTorrent StepDelivery = "waiting_user_chose_torrent"
	// GetMagnetLink получение магнет ссылки
	GetMagnetLink StepDelivery = "get_magnet_link_status"
	// AddTorrentToTorrentClient Добавление раздачи для скачивания торрент клиентом
	AddTorrentToTorrentClient StepDelivery = "add_torrent_client_status"
	// PrepareFileMatches получение информации о файлах раздачи
	PrepareFileMatches StepDelivery = "prepare_file_matches"
	// WaitingChoseFileMatches ожидание подтверждения пользователем соответствий выбора файлов
	WaitingChoseFileMatches StepDelivery = "waiting_chose_file_matches"
	// WaitingTorrentDownloadComplete ожидание завершения окончания скачивания раздачи
	WaitingTorrentDownloadComplete StepDelivery = "waiting_torrent_download_complete"
	// CreateVideoContentCatalogs Формирование каталогов и иерархии файлов
	CreateVideoContentCatalogs StepDelivery = "create_video_content_catalogs"
	// DeterminingNeedConvertFiles Определение необходимости конвертации файлов
	DeterminingNeedConvertFiles StepDelivery = "determining_need_convert_files"
	// --- Ветвь если необходимо добавление аудио дорожек/субтитров

	// MergeVideoFiles Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
	MergeVideoFiles StepDelivery = "merge_video_files"
	// -- Ветвь если не нужно изменять исходные файлы

	// CopyVideoFiles Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
	CopyVideoFiles StepDelivery = "copy_video_files"
	// SetMediaMetaData установка методаных серий сезона сериала / фильма в медиасервере
	SetMediaMetaData StepDelivery = "set_media_meta_data"
	// SendDeliveryNotification Отправка уведомления в telegramm о успешной доставки видеофайлов до медиа сервера
	SendDeliveryNotification StepDelivery = "send_delivery_notification"
)

// ContentDeliveryData модель содержащая информацию о видео контенте для сезона сериала / фильма
/*
	К одному сезону сериала / фильму может быть привязано несколько VideoContent,
	но для упрощения пока будем пока разрешать только 1
*/
type ContentDeliveryData struct {
	// Данные выпуска
	SearchQuery           *string
	TorrentSearch         *contentdelivery.TorrentSearchResult
	SelectTorrentHref     *string
	MagnetInfo            *contentdelivery.MagnetInfo
	ContentMatches        []contentdelivery.ContentMatches
	TorrentDownloadStatus *contentdelivery.TorrentDownloadStatus
	CatalogsInfo          *contentdelivery.CatalogsInfo
	MergeVideoStatus      *contentdelivery.MergeVideoStatus
}

type ContentDeliveryMetadata struct {
	MediaID contentdelivery.MediaID
}

type ContentDeliveryFailData struct {
}

type CreateOptions struct {
	MediaID contentdelivery.MediaID
}

func (c CreateOptions) GetIdempotencyKey() string {

	if c.MediaID.MovieID != nil {
		return fmt.Sprintf("movie_%d", *c.MediaID.MovieID)
	}

	if c.MediaID.TVShow != nil {
		return fmt.Sprintf("tv_%d_season_%d", c.MediaID.TVShow.TVShowID, c.MediaID.TVShow.SeasonNumber)
	}

	return ""
}
