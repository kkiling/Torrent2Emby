package config

import (
	"github.com/joho/godotenv"

	"github.com/kkiling/torrent2emby/internal/log"
)

// Константы для имен переменных окружения
const (
	TheMovieDbApiKey     = "THE_MOVIE_DB_API_KEY"
	RutrackerUsername    = "RUTRACKER_USERNAME"
	RutrackerPassword    = "RUTRACKER_PASSWORD"
	RutrackerCookieDir   = "RUTRACKER_COOKIE_DIR"
	QBittorrentUsername  = "QBITTORRENT_USERNAME"
	QBittorrentPassword  = "QBITTORRENT_PASSWORD"
	QBittorrentCookieDir = "QBITTORRENT_COOKIE_DIR"
	QBittorrentApiUrl    = "QBITTORRENT_API_URL"
	EmbyApiUrl           = "EMBY_API_URL"
	EmbyApiKey           = "EMBY_API_KEY"
	SqliteDsn            = "SQLITE_DSN"
)

// MovieDbConfig конфигурация для The Movie DB API
type MovieDbConfig struct {
	ApiKey string
}

// RutrackerConfig конфигурация для Rutracker
type RutrackerConfig struct {
	Username  string
	Password  string
	CookieDir string
}

// QBittorrentConfig конфигурация для QBittorrent
type QBittorrentConfig struct {
	Username  string
	Password  string
	CookieDir string
	ApiUrl    string
}

// EmbyConfig конфигурация для Emby Api
type EmbyConfig struct {
	ApiKey string
	ApiUrl string
}

type StorageConfig struct {
	SqliteDsn string
}

// EnvConfig объединяет все конфигурации
type EnvConfig struct {
	MovieDb     MovieDbConfig
	Rutracker   RutrackerConfig
	QBittorrent QBittorrentConfig
	Emby        EmbyConfig
	Storage     StorageConfig
}

func loadMovieDbConfig() (*MovieDbConfig, error) {
	apiKey, err := getEnvString(TheMovieDbApiKey)
	if err != nil {
		return nil, err
	}

	return &MovieDbConfig{
		ApiKey: apiKey,
	}, nil
}

func loadRutrackerConfig() (*RutrackerConfig, error) {
	username, err := getEnvString(RutrackerUsername)
	if err != nil {
		return nil, err
	}

	password, err := getEnvString(RutrackerPassword)
	if err != nil {
		return nil, err
	}

	cookieDir, err := getEnvString(RutrackerCookieDir)
	if err != nil {
		return nil, err
	}

	return &RutrackerConfig{
		Username:  username,
		Password:  password,
		CookieDir: cookieDir,
	}, nil
}

func loadEmbyConfig() (*EmbyConfig, error) {
	apiKey, err := getEnvString(EmbyApiKey)
	if err != nil {
		return nil, err
	}

	url, err := getEnvString(EmbyApiUrl)
	if err != nil {
		return nil, err
	}

	return &EmbyConfig{
		ApiKey: apiKey,
		ApiUrl: url,
	}, nil
}

func loadQBittorrentConfig() (*QBittorrentConfig, error) {
	username, err := getEnvString(QBittorrentUsername)
	if err != nil {
		return nil, err
	}

	password, err := getEnvString(QBittorrentPassword)
	if err != nil {
		return nil, err
	}

	cookieDir, err := getEnvString(QBittorrentCookieDir)
	if err != nil {
		return nil, err
	}

	apiUrl, err := getEnvString(QBittorrentApiUrl)
	if err != nil {
		return nil, err
	}

	return &QBittorrentConfig{
		Username:  username,
		Password:  password,
		CookieDir: cookieDir,
		ApiUrl:    apiUrl,
	}, nil
}

func loadStorageConfig() (*StorageConfig, error) {
	SqliteDsn, err := getEnvString(SqliteDsn)
	if err != nil {
		return nil, err
	}

	return &StorageConfig{
		SqliteDsn: SqliteDsn,
	}, nil
}

func NewEnvConfig(logger log.Logger) (*EnvConfig, error) {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Warning: Could not find .env prepare - using system environment variables")
	}

	movieDbConfig, err := loadMovieDbConfig()
	if err != nil {
		return nil, err
	}

	rutrackerConfig, err := loadRutrackerConfig()
	if err != nil {
		return nil, err
	}

	qBittorrentConfig, err := loadQBittorrentConfig()
	if err != nil {
		return nil, err
	}

	embyConfig, err := loadEmbyConfig()
	if err != nil {
		return nil, err
	}

	storageConfig, err := loadStorageConfig()
	if err != nil {
		return nil, err
	}

	return &EnvConfig{
		MovieDb:     *movieDbConfig,
		Rutracker:   *rutrackerConfig,
		QBittorrent: *qBittorrentConfig,
		Emby:        *embyConfig,
		Storage:     *storageConfig,
	}, nil
}
