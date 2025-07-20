package runner

import (
	"context"

	"github.com/kkiling/torrent2emby/internal/usercase/videodelivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params videodelivery.GenerateSearchQueryParams) (string, error)
	SearchTorrent(ctx context.Context, params videodelivery.SearchTorrentParams) (*videodelivery.TorrentSearchResult, error)
	GetMagnetLink(ctx context.Context, params videodelivery.GetMagnetLinkParams) (*videodelivery.MagnetInfo, error)
	AddTorrentToTorrentClient(ctx context.Context, params videodelivery.AddTorrentParams) error
	PrepareFileMatches(ctx context.Context, params videodelivery.PreparingFileMatchesParams) ([]videodelivery.ContentMatches, error)
}
