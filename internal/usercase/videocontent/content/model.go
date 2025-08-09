package content

import (
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners"
)

type DeliveryStatus string

const (
	// DeliveryStatusFailed доставка была зафейлена
	DeliveryStatusFailed DeliveryStatus = "failed"
	// DeliveryStatusInProgress - В процессе доставки файлов
	DeliveryStatusInProgress DeliveryStatus = "in_progress"
	// DeliveryStatusDelivered - Файлы доставлены
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	// DeliveryStatusUpdating - Обновление раздачи
	DeliveryStatusUpdating DeliveryStatus = "updating"
	// DeliveryStatusDeleting - В процессе удаления файлов в диска
	DeliveryStatusDeleting DeliveryStatus = "deleting"
	// DeliveryStatusDeleted - Файлы удалены
	DeliveryStatusDeleted DeliveryStatus = "deleted"
)

type TorrentInfo struct {
	// href ссылки на торрент сайт раздачи
	Href string
	// Magnet ссылка текущую раздачу
	Magnet string
	// Хеш торрента
	Hash string
}

// VideoContent информация о файлах
type VideoContent struct {
	// ID информации о файлах
	ID uuid.UUID
	// CreatedAt Датоа создания
	CreatedAt time.Time
	// ID сериала/фильма
	ContentID common.ContentID
	// Инфа о торренте
	TorrentInfo *TorrentInfo
	// Статус
	DeliveryStatus DeliveryStatus
	// Стейты привязанные к текущему контенту
	State []State
}

// State Таблица выпусков связанных с TVShowContent
type State struct {
	StateID uuid.UUID
	Type    runners.Type
}

type UpdateVideoContent struct {
	// Инфа о торренте
	TorrentInfo *TorrentInfo
	// Статус
	DeliveryStatus DeliveryStatus
}
