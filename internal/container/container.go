package container

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/emby"
	"github.com/kkiling/torrent2emby/internal/runner/tvshowdeliverystate"

	prepareTVShow "github.com/kkiling/torrent2emby/internal/adapter/matchtvshow"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	mkvsqlite "github.com/kkiling/torrent2emby/internal/adapter/mkvmerge/storage/sqlite"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/adapter/rutracker"
	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	statemachinesqlite "github.com/kkiling/torrent2emby/internal/statemachine/storage/sqlite"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowdelivery"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary/storage/sqlite"
)

type Container struct {
	tvShowLibrary        *tvshowlibrary.Service
	deliveryStateMachine *tvshowdeliverystate.StateMachineService
	mkvMergePipeline     *mkvmerge.Pipeline
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

	mkvPipelineStorage, err := mkvsqlite.NewStorage(mkvsqlite.Config{
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

	embyApi, err := emby.NewApi(cfg.Emby.ApiKey, cfg.Emby.ApiUrl, logger)
	if err != nil {
		return nil, fmt.Errorf("emby.NewApi: %w", err)
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

	mkvMerge := mkvmerge.NewMerge(logger)
	mkvPipeline := mkvmerge.NewPipeline(mkvMerge, mkvPipelineStorage, logger)

	prepareTVShowService := prepareTVShow.NewService(mkvMerge)

	// UserCase
	tvShowLibrary := tvshowlibrary.NewService(tvShowLibraryStorage, themoviedbApi)

	delivery := tvshowdelivery.NewService(
		// TODO: вынести в конфиг
		tvshowdelivery.Config{
			BasePath:                   "/nfs",
			TVShowTorrentSavePath:      "/downloads",
			TvShowMediaSaveTvShowsPath: "/tvshows",
			UserGroup:                  "nas",
		},
		tvShowLibrary,
		rutrackerApi,
		qBittorrentApi,
		embyApi,
		prepareTVShowService,
		mkvPipeline,
	)
	deliveryStateMachine := tvshowdeliverystate.NewState(delivery, stateStorage)

	return &Container{
		tvShowLibrary:        tvShowLibrary,
		deliveryStateMachine: deliveryStateMachine,
		mkvMergePipeline:     mkvPipeline,
	}, nil
}

func (c *Container) GetTvShowLibrary() *tvshowlibrary.Service {
	return c.tvShowLibrary
}

func (c *Container) DeliveryStateMachine() *tvshowdeliverystate.StateMachineService {
	return c.deliveryStateMachine
}

func (c *Container) MkvMergePipeline() *mkvmerge.Pipeline {
	return c.mkvMergePipeline
}
