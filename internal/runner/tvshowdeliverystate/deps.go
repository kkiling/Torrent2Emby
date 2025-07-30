package tvshowdeliverystate

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/usercase/tvshowdelivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params tvshowdelivery.GenerateSearchQueryParams) (string, error)
	SearchTorrent(ctx context.Context, params tvshowdelivery.SearchTorrentParams) (*tvshowdelivery.TorrentSearchResult, error)
	GetMagnetLink(ctx context.Context, params tvshowdelivery.GetMagnetLinkParams) (*tvshowdelivery.MagnetInfo, error)
	AddTorrentToTorrentClient(ctx context.Context, params tvshowdelivery.AddTorrentParams) error
	PrepareFileMatches(ctx context.Context, params tvshowdelivery.PreparingFileMatchesParams) ([]tvshowdelivery.ContentMatches, error)
	WaitingTorrentDownloadComplete(ctx context.Context, params tvshowdelivery.WaitingTorrentDownloadCompleteParams) (*tvshowdelivery.TorrentDownloadStatus, error)
	CreateContentCatalogs(ctx context.Context, params tvshowdelivery.CreateContentCatalogsParams) (*tvshowdelivery.CatalogsInfo, error)
	StartMergeVideo(ctx context.Context, params tvshowdelivery.MergeVideoParams) ([]tvshowdelivery.MergeVideoFile, error)
	GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*tvshowdelivery.MergeVideoStatus, error)
	SetVideoFileGroup(ctx context.Context, files []string) error
	SetMediaMetaData(ctx context.Context, params tvshowdelivery.SetMediaMetaDataParams) error
}
