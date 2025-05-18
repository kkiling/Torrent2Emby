from enum import Enum
from typing import Optional, List
from dataclasses import dataclass
from datetime import datetime
from typing import cast

@dataclass
class TorrentAddOptions:
    urls: str
    save_path: str
    category: Optional[str] = None
    tags: Optional[List[str]] = None
    paused: bool = False


@dataclass
class TorrentState(Enum):
    """Состояния торрента"""
    ERROR = "error"  # Произошла ошибка, применяется к приостановленным торрентам
    MISSING_FILES = "missingFiles"  # Файлы данных торрента отсутствуют
    UPLOADING = "uploading"  # Торрент раздается и передает данные
    PAUSED_UP = "pausedUP"  # Торрент приостановлен и завершил загрузку
    QUEUED_UP = "queuedUP"  # Очередь включена и торрент в очереди на раздачу
    STALLED_UP = "stalledUP"  # Торрент раздается, но нет активных соединений
    CHECKING_UP = "checkingUP"  # Торрент завершил загрузку и проверяется
    FORCED_UP = "forcedUP"  # Торрент принудительно раздается, игнорируя ограничение очереди
    ALLOCATING = "allocating"  # Выделение дискового пространства для загрузки
    DOWNLOADING = "downloading"  # Торрент загружается и передает данные
    META_DL = "metaDL"  # Торрент только начал загрузку и получает метаданные
    PAUSED_DL = "pausedDL"  # Торрент приостановлен и НЕ завершил загрузку
    QUEUED_DL = "queuedDL"  # Очередь включена и торрент в очереди на загрузку
    STALLED_DL = "stalledDL"  # Торрент загружается, но нет активных соединений
    CHECKING_DL = "checkingDL"  # Проверка, но торрент НЕ завершил загрузку
    FORCED_DL = "forcedDL"  # Торрент принудительно загружается, игнорируя очередь
    CHECKING_RESUME_DATA = "checkingResumeData"  # Проверка данных возобновления при запуске
    MOVING = "moving"  # Торрент перемещается в другое место
    UNKNOWN = "unknown"  # Неизвестный статус

@dataclass
class TorrentInfo:
    hash: str  # Хеш торрента
    name: str  # Имя торрента
    category: str  # Категория торрента
    tags: str  # Список тегов через запятую
    content_path: str  # Абсолютный путь к содержимому торрента
    state: TorrentState  # Состояние торрента
    added_on: datetime  # Время добавления торрента (Unix Epoch)
    completion_on: datetime  # Время завершения загрузки (Unix Epoch)
    eta: int  # Оставшееся время (в секундах)
    amount_left: int  # Количество оставшихся для загрузки данных (байт)
    completed: int  # Количество загруженных данных (байт)
    downloaded: int  # Общее количество загруженных данных (байт)
    uploaded: int  # Количество отданных данных (байт)
    size: int  # Размер выбранных файлов (байт)
    total_size: int  # Общий размер всех файлов (байт)
    progress: float  # Прогресс торрента (процент/100)
    dl_speed: int  # Текущая скорость загрузки (байт/с)
    up_speed: int  # Текущая скорость отдачи (байт/с)
