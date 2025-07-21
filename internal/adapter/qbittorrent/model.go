package qbittorrent

import "time"

// TorrentState представляет состояние торрента
type TorrentState string

const (
	// TorrentStateError Произошла ошибка, применяется к приостановленным торрентам
	TorrentStateError TorrentState = "error"
	// TorrentStateMissingFiles Файлы данных торрента отсутствуют
	TorrentStateMissingFiles TorrentState = "missingFiles"
	// TorrentStateUploading Торрент раздается и данные передаются
	TorrentStateUploading TorrentState = "uploading"
	// TorrentStatePausedUP Торрент приостановлен и загрузка завершена
	TorrentStatePausedUP TorrentState = "pausedUP"
	// TorrentStateQueuedUP Очередь включена и торрент в очереди на раздачу
	TorrentStateQueuedUP TorrentState = "queuedUP"
	// TorrentStateStalledUP Торрент раздается, но нет активных соединений
	TorrentStateStalledUP TorrentState = "stalledUP"
	// TorrentStateCheckingUP Торрент завершил загрузку и производится проверка
	TorrentStateCheckingUP TorrentState = "checkingUP"
	// TorrentStateForcedUP Торрент принудительно раздается, игнорируя ограничение очереди
	TorrentStateForcedUP TorrentState = "forcedUP"
	// TorrentStateAllocating Торрент выделяет дисковое пространство для загрузки
	TorrentStateAllocating TorrentState = "allocating"
	// TorrentStateDownloading Торрент загружается и данные передаются
	TorrentStateDownloading TorrentState = "downloading"
	// TorrentStateMetaDL Торрент только начал загрузку и получает метаданные
	TorrentStateMetaDL TorrentState = "metaDL"
	// TorrentStatePausedDL Торрент приостановлен и загрузка НЕ завершена
	TorrentStatePausedDL TorrentState = "pausedDL"
	// TorrentStateStoppedDL Торрент приостановлен и загрузка НЕ завершена
	TorrentStateStoppedDL TorrentState = "stoppedDL"
	// TorrentStateQueuedDL Очередь включена и торрент в очереди на загрузку
	TorrentStateQueuedDL TorrentState = "queuedDL"
	// TorrentStateStalledDL Торрент загружается, но нет активных соединений
	TorrentStateStalledDL TorrentState = "stalledDL"
	// TorrentStateCheckingDL Аналогично checkingUP, но торрент НЕ завершил загрузку
	TorrentStateCheckingDL TorrentState = "checkingDL"
	// TorrentStateForcedDL Торрент принудительно загружается, игнорируя ограничение очереди
	TorrentStateForcedDL TorrentState = "forcedDL"
	// TorrentStateCheckingResumeData Проверка данных возобновления при запуске qBittorrent
	TorrentStateCheckingResumeData TorrentState = "checkingResumeData"
	// TorrentStateMoving Торрент перемещается в другое место
	TorrentStateMoving TorrentState = "moving"
	// TorrentStateUnknown Неизвестный статус
	TorrentStateUnknown TorrentState = "unknown"
)

// TorrentInfo содержит информацию о торренте
type TorrentInfo struct {
	// Хеш торрента
	Hash string
	// Название торрента
	Name string
	// Категория торрента
	Category string
	// Теги торрента, разделенные запятой
	Tags string
	// Путь сохранения торрента
	SavePath string
	// Абсолютный путь к содержимому торрента (корневой путь для торрентов с несколькими файлами,
	// абсолютный путь к файлу для торрентов с одним файлом)
	ContentPath string
	// Текущее состояние торрента
	State TorrentState
	// Время добавления торрента в клиент
	AddedOn time.Time
	// Время завершения загрузки торрента
	CompletionOn time.Time
	// Расчетное время до завершения загрузки (в секундах)
	Eta int
	// Количество оставшихся для загрузки данных (в байтах)
	AmountLeft int64
	// Количество завершенных данных (в байтах)
	Completed int64
	// Общее количество загруженных данных (в байтах)
	Downloaded int64
	// Общее количество отданных данных (в байтах)
	Uploaded int64
	// Размер выбранных для загрузки файлов (в байтах)
	Size int64
	// Общий размер всех файлов в торренте, включая невыбранные (в байтах)
	TotalSize int64
	// Прогресс загрузки торрента (от 0 до 1)
	Progress float64
	// Текущая скорость загрузки (байт/сек)
	DlSpeed int64
	// Текущая скорость отдачи (байт/сек)
	UpSpeed int64
}

// TorrentFile содержит информацию о файле в торренте
type TorrentFile struct {
	Index int
	// Имя файла с полным путем внутри торрента
	Name string
	// Прогресс загрузки файла (от 0 до 1)
	Progress float64
	// Размер файла в байтах
	Size int64
}
