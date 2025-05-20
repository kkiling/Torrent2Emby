package config

import (
	"github.com/joho/godotenv"
	"github.com/kkiling/torrent2emby/internal/log"
)

// Константы для имен переменных окружения
const (
	TheMovieDbApiKey = "THE_MOVIE_DB_API_KEY"

	RutrackerUsername  = "RUTRACKER_USERNAME"
	RutrackerPassword  = "RUTRACKER_PASSWORD"
	RutrackerCookieDir = "RUTRACKER_COOKIE_DIR"

	QBittorrentUsername  = "QBITTORRENT_USERNAME"
	QBittorrentPassword  = "QBITTORRENT_PASSWORD"
	QBittorrentCookieDir = "QBITTORRENT_COOKIE_DIR"
	QBittorrentApiUrl    = "QBITTORRENT_API_URL"
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

// EnvConfig объединяет все конфигурации
type EnvConfig struct {
	MovieDb     MovieDbConfig
	Rutracker   RutrackerConfig
	QBittorrent QBittorrentConfig
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

func NewEnvConfig(logger log.Logger) (*EnvConfig, error) {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Warning: Could not find .env file - using system environment variables")
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

	return &EnvConfig{
		MovieDb:     *movieDbConfig,
		Rutracker:   *rutrackerConfig,
		QBittorrent: *qBittorrentConfig,
	}, nil
}
