package server

import (
	"context"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"

	"github.com/kkiling/torrent2emby/internal/handler/tvshowlibrary"
)

// Torrent2EmbyServer сервер
type Torrent2EmbyServer struct {
	*CustomServer
}

// NewTorrent2EmbyServer новый сервер
func NewTorrent2EmbyServer(
	logger log.Logger,
	cfg server.Config,
	tvShowLibrary tvshowlibrary.TVShowLibrary,
) *Torrent2EmbyServer {
	return &Torrent2EmbyServer{
		CustomServer: NewCustomServer(
			logger,
			cfg,
			tvshowlibrary.NewHandler(logger, tvShowLibrary),
		),
	}
}

// Start старт сервера
func (s *Torrent2EmbyServer) Start(ctx context.Context) error {
	return s.CustomServer.Start(ctx, "api")
}
