package config

import (
	"fmt"

	"github.com/kkiling/goplatform/config"
)

const (
	// ServerConfigName конфиг сервера
	ServerConfigName = "server"
	// TheMovieDbName конфиг TheMovieDb api
	TheMovieDbName = "the_movie_db"
	// RutrackerName конфиг Rutracker api
	RutrackerName = "rutracker"
	// QBitTorrentName конфиг QBitTorrent api
	QBitTorrentName = "qbittorrent"
	// EmbyName конфиг Emby api
	EmbyName = "emby"
	// SqliteName конфиг sqlite
	SqliteName = "sqlite"
	// DeliveryName конфиг sqlite
	DeliveryName = "delivery"
)

// AppConfig объединяет все конфигурации
type AppConfig struct {
	Server         ServerConfig      `yaml:"server"`
	TheMovieDb     TheMovieDbConfig  `yaml:"the_movie_db"`
	Rutracker      RutrackerConfig   `yaml:"rutracker"`
	QBittorrent    QBittorrentConfig `yaml:"qbittorrent"`
	Emby           EmbyConfig        `yaml:"emby"`
	Sqlite         SqliteConfgi      `yaml:"sqlite"`
	DeliveryConfig DeliveryConfig    `yaml:"delivery"`
}

// TheMovieDbConfig конфигурация для The Movie DB API
type TheMovieDbConfig struct {
	ApiKey string `yaml:"api_key"`
}

// RutrackerConfig конфигурация для Rutracker
type RutrackerConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	CookieDir string `yaml:"cookie_dir"`
}

// QBittorrentConfig конфигурация для QBittorrent
type QBittorrentConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	CookieDir string `yaml:"cookie_dir"`
	ApiUrl    string `yaml:"api_url"`
}

// EmbyConfig конфигурация для Emby Api
type EmbyConfig struct {
	ApiKey string `yaml:"api_key"`
	ApiUrl string `yaml:"api_url"`
}

type SqliteConfgi struct {
	SqliteDsn string `yaml:"sqlite_dsn"`
}

type ServerConfig struct {
	Host                    string `yaml:"host"`
	GrpcPort                int    `yaml:"grpc_port"`
	HttpPort                int    `yaml:"http_port"`
	MaxSendMessageLength    int    `yaml:"max_send_message_length"`
	MaxReceiveMessageLength int    `yaml:"max_receive_message_length"`
	ShutdownTimeout         int    `yaml:"shutdown_timeout"`
}

// DeliveryConfig конфиг для доставки контента
type DeliveryConfig struct {
	BasePath                   string `yaml:"base_path"`
	TVShowTorrentSavePath      string `yaml:"tv_show_torrent_save_path"`
	TVShowMediaSaveTvShowsPath string `yaml:"tv_show_media_save_tv_shows_path"`
	UserGroup                  string `yaml:"user_group"`
}

func loadCfg[T any](cfgName string, cfgProvider config.Provider) (*T, error) {
	var cfg T
	err := cfgProvider.PopulateByKey(cfgName, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cfgProvider.PopulateByKey: %w", err)
	}
	return &cfg, nil
}

func NewEnvConfig(cfgProvider config.Provider) (*AppConfig, error) {
	movieDbConfig, err := loadCfg[TheMovieDbConfig](TheMovieDbName, cfgProvider)
	if err != nil {
		return nil, err
	}

	rutrackerConfig, err := loadCfg[RutrackerConfig](RutrackerName, cfgProvider)
	if err != nil {
		return nil, err
	}

	qBittorrentConfig, err := loadCfg[QBittorrentConfig](QBitTorrentName, cfgProvider)
	if err != nil {
		return nil, err
	}

	embyConfig, err := loadCfg[EmbyConfig](EmbyName, cfgProvider)
	if err != nil {
		return nil, err
	}

	sqliteConfgi, err := loadCfg[SqliteConfgi](SqliteName, cfgProvider)
	if err != nil {
		return nil, err
	}

	serverConfig, err := loadCfg[ServerConfig](ServerConfigName, cfgProvider)
	if err != nil {
		return nil, err
	}

	deliveryConfig, err := loadCfg[DeliveryConfig](DeliveryName, cfgProvider)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Server:         *serverConfig,
		TheMovieDb:     *movieDbConfig,
		Rutracker:      *rutrackerConfig,
		QBittorrent:    *qBittorrentConfig,
		Emby:           *embyConfig,
		Sqlite:         *sqliteConfgi,
		DeliveryConfig: *deliveryConfig,
	}, nil
}
