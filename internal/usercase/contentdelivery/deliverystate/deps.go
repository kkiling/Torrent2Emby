package deliverystate

import (
	"context"

	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params contentdelivery.GenerateSearchQueryParams) (string, error)
	SearchTorrent(ctx context.Context, params contentdelivery.SearchTorrentParams) (*contentdelivery.TorrentSearchResult, error)
	GetMagnetLink(ctx context.Context, params contentdelivery.GetMagnetLinkParams) (*contentdelivery.MagnetInfo, error)
	AddTorrentToTorrentClient(ctx context.Context, params contentdelivery.AddTorrentParams) error
	PrepareFileMatches(ctx context.Context, params contentdelivery.PreparingFileMatchesParams) ([]contentdelivery.ContentMatches, error)
}
