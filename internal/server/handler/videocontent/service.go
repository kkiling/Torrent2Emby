package videocontent

import (
	"context"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/content"
	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
	"google.golang.org/grpc"
	"net/http"
)

// VideoContent юзеркейс работы с доставкой видео файлов
type VideoContent interface {
	CreateVideoContent(ctx context.Context, params content.CreateVideoContentParams) (*videocontent.VideoContent, error)
	GetVideoContent(ctx context.Context, contentID videocontent.ContentID) ([]videocontent.VideoContent, error)
	GetTVShowDeliveryData(ctx context.Context, contentID videocontent.ContentID) (*videocontent.TVShowDeliveryState, error)
	ChoseTorrentOptions(ctx context.Context, contentID videocontent.ContentID, opts videocontent.ChoseTorrentOptions) (*videocontent.TVShowDeliveryState, error)
	ChoseFileMatchesOptions(ctx context.Context, contentID videocontent.ContentID, opts videocontent.ChoseFileMatchesOptions) (*videocontent.TVShowDeliveryState, error)
}

type Handler struct {
	desc.VideoContentServiceServer
	logger       log.Logger
	videoContent VideoContent
}

// NewHandler новый хендлер
func NewHandler(logger log.Logger,
	videoContent VideoContent,
) *Handler {
	return &Handler{
		logger:       logger.Named("video_content"),
		videoContent: videoContent,
	}
}

// RegistrationServerHandlers .
func (h *Handler) RegistrationServerHandlers(_ *http.ServeMux) {
}

// RegisterServiceHandlerFromEndpoint .
func (h *Handler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterVideoContentServiceHandlerFromEndpoint
}

// RegisterServiceServer регистрация
func (h *Handler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterVideoContentServiceServer(server, h)
}
