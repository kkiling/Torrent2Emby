package content

import (
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent"
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners"
)

type DeliveryStatus string

const (
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

// VideoContent информация о файлах
type VideoContent struct {
	// ID информации о файлах
	ID uuid.UUID
	// CreatedAt Датоа создания
	CreatedAt time.Time
	// ID сериала/фильма
	ContentID videocontent.ContentID
	// href ссылки на торрент сайт раздачи
	Href *string
	// Magnet ссылка текущую раздачу
	Magnet *string
	// Хеш торрента
	Hash *string
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

/*
	Какие кейсы:
		- Создание новой доставки файлов с прохождением полного флоу от поиска раздачи до доставки файлов до медиасервера
		- Можно оставить информацию о раздаче (Href и Magnet) но при этом удалить все файлы, что бы не занимали место на диске
		- Потом на основе (Href и Magnet) восстанавливать файлы и скачивать их снова, при этом не запрашивая больше инфу от клиента
			и все подтягивать из старых стейтов (что делать если раздача обновиться?)
       - Раздача может обновиться и запускается процесс обновления раздачи
*/
