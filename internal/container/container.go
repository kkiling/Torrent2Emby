package container

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	prepareTVShow "github.com/kkiling/torrent2emby/internal/adapter/prepare/tvshow"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/adapter/rutracker"
	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	statemachinesqlite "github.com/kkiling/torrent2emby/internal/statemachine/storage/sqlite"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery/deliverystate"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary/storage/sqlite"
)

type Container struct {
	tvShowLibrary        *tvshowlibrary.Service
	deliveryStateMachine *deliverystate.StateMachineService
}

func NewContainer() (*Container, error) {
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		return nil, fmt.Errorf("config.NewEnvConfig: %w", err)
	}

	// Storage
	tvShowLibraryStorage, err := sqlite.NewStorage(sqlite.Config{
		DSN: cfg.Storage.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	stateStorage, err := statemachinesqlite.NewStorage(statemachinesqlite.Config{
		DSN: cfg.Storage.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	// Adapter
	themoviedbApi, err := themoviedb.NewApi(
		cfg.MovieDb.ApiKey,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("themoviedb.NewApi: %w", err)
	}

	qBittorrentApi, err := qbittorrent.NewApi(
		logger,
		cfg.QBittorrent.ApiUrl,
		cfg.QBittorrent.Username,
		cfg.QBittorrent.Password,
		cfg.QBittorrent.CookieDir,
	)
	if err != nil {
		return nil, fmt.Errorf("qbittorrent.NewApi: %w", err)
	}

	rutrackerApi, err := rutracker.NewApi(
		logger,
		cfg.Rutracker.Username,
		cfg.Rutracker.Password,
		cfg.Rutracker.CookieDir,
	)
	if err != nil {
		return nil, fmt.Errorf("rutracker.NewApi: %w", err)
	}

	mkvMerge := mkvmerge.NewService()
	prepareTVShow := prepareTVShow.NewService(mkvMerge)

	// UserCase
	tvShowLibrary := tvshowlibrary.NewService(tvShowLibraryStorage, themoviedbApi)

	delivery := contentdelivery.NewService(
		// TODO: вынести в конфиг
		contentdelivery.Config{
			BasePath:                   "/nfs",
			TVShowTorrentSavePath:      "/downloads",
			TvShowMediaSaveTvShowsPath: "/tvshows",
			UserGroup:                  "nas",
		},
		tvShowLibrary,
		rutrackerApi,
		qBittorrentApi,
		prepareTVShow,
		mkvMerge,
	)
	deliveryStateMachine := deliverystate.NewState(delivery, stateStorage)

	return &Container{
		tvShowLibrary:        tvShowLibrary,
		deliveryStateMachine: deliveryStateMachine,
	}, nil
}

func (c *Container) GetTvShowLibrary() *tvshowlibrary.Service {
	return c.tvShowLibrary
}

func (c *Container) DeliveryStateMachine() *deliverystate.StateMachineService {
	return c.deliveryStateMachine
}
