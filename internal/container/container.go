package container

import (
	"fmt"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase/sqlitebase"
	"github.com/kkiling/statemachine"

	"github.com/kkiling/torrent2emby/internal/adapter/emby"
	prepareTVShow "github.com/kkiling/torrent2emby/internal/adapter/matchtvshow"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	mkvsqlite "github.com/kkiling/torrent2emby/internal/adapter/mkvmerge/storage/sqlite"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/adapter/rutracker"
	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary/storage/sqlite"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/delivery"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type Container struct {
	logger                     log.Logger
	tvShowLibrary              *tvshowlibrary.Service
	tvShowDeliveryStateMachine *tvshowdeliverystate.StateMachineService
	mkvMergePipeline           *mkvmerge.Pipeline
}

func NewContainer(cfg *config.AppConfig) (*Container, error) {
	logger := log.NewLogger(log.DebugLevel)

	// Storage
	tvShowLibraryStorage, err := sqlite.NewStorage(sqlitebase.Config{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	stateStorage, err := statemachine.NewSqliteStorage(statemachine.SqliteConfig{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	mkvPipelineStorage, err := mkvsqlite.NewStorage(sqlitebase.Config{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	// Adapter
	themoviedbApi, err := themoviedb.NewApi(
		cfg.TheMovieDb.ApiKey,
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

	deliveryService := delivery.NewService(
		delivery.Config{
			BasePath:                   cfg.DeliveryConfig.BasePath,
			TVShowTorrentSavePath:      cfg.DeliveryConfig.TVShowTorrentSavePath,
			TVShowMediaSaveTvShowsPath: cfg.DeliveryConfig.TVShowMediaSaveTvShowsPath,
			UserGroup:                  cfg.DeliveryConfig.UserGroup,
		},
		tvShowLibrary,
		rutrackerApi,
		qBittorrentApi,
		embyApi,
		prepareTVShowService,
		mkvPipeline,
	)
	tvShowDeliveryStateMachine := tvshowdeliverystate.NewState(deliveryService, stateStorage)

	return &Container{
		logger:                     logger,
		tvShowLibrary:              tvShowLibrary,
		tvShowDeliveryStateMachine: tvShowDeliveryStateMachine,
		mkvMergePipeline:           mkvPipeline,
	}, nil
}

func (c *Container) GetTvShowLibrary() *tvshowlibrary.Service {
	return c.tvShowLibrary
}

func (c *Container) TVShowDeliveryStateMachine() *tvshowdeliverystate.StateMachineService {
	return c.tvShowDeliveryStateMachine
}

func (c *Container) MkvMergePipeline() *mkvmerge.Pipeline {
	return c.mkvMergePipeline
}

func (c *Container) GetLogger() log.Logger {
	return c.logger
}
