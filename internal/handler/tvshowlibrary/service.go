package tvshowlibrary

import (
	"net/http"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"
	"google.golang.org/grpc"

	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
)

// TVShowLibrary юзеркейс авторизации
type TVShowLibrary interface {
}

type Handler struct {
	desc.TVShowLibraryServiceServer
	logger        log.Logger
	tvShowLibrary TVShowLibrary
}

// NewHandler новый хендлер
func NewHandler(logger log.Logger,
	tvShowLibrary TVShowLibrary,
) *Handler {
	return &Handler{
		logger:        logger.Named("tv_show_library"),
		tvShowLibrary: tvShowLibrary,
	}
}

// RegistrationServerHandlers .
func (h *Handler) RegistrationServerHandlers(_ *http.ServeMux) {
}

// RegisterServiceHandlerFromEndpoint .
func (h *Handler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterTVShowLibraryServiceHandlerFromEndpoint
}

// RegisterServiceServer регистрация
func (h *Handler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterTVShowLibraryServiceServer(server, h)
}
