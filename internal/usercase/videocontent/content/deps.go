package content

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners/tvshowdeliverystate"
	"time"
)

type Storage interface {
	SaveVideoContent(ctx context.Context, videoContent *VideoContent) error
	GetVideoContents(ctx context.Context, contentID common.ContentID) ([]VideoContent, error)
	UpdateVideoContent(ctx context.Context, id uuid.UUID, videoContent *UpdateVideoContent) error
	GetVideoContentsByStatus(ctx context.Context, status DeliveryStatus, limit int) ([]VideoContent, error)
}

type TVShowLibrary interface {
	GetTVShowInfo(ctx context.Context, params tvshowlibrary.GetTVShowParams) (*tvshowlibrary.GetTVShowResult, error)
	GetSeasonEpisodes(ctx context.Context, params tvshowlibrary.GetSeasonEpisodesParams) (*tvshowlibrary.GetSeasonEpisodesResult, error)
}

type TVShowDeliveryState interface {
	GetStateByID(ctx context.Context, stateID uuid.UUID) (*tvshowdeliverystate.State, error)
	Create(ctx context.Context, opt tvshowdeliverystate.CreateOptions) (*tvshowdeliverystate.State, error)
	Complete(ctx context.Context, stateID uuid.UUID, options ...any) (st *tvshowdeliverystate.State, executeErr error, err error)
}

// UUIDGenerator интерфейс для генерации UUID (реальный или мок)
type UUIDGenerator interface {
	New() uuid.UUID
}

// Clock интерфейс для работы со временем (реальный или мок)
type Clock interface {
	Now() time.Time
}

type uuidGenerator struct{}

func (uuidGenerator) New() uuid.UUID {
	return uuid.New()
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}
