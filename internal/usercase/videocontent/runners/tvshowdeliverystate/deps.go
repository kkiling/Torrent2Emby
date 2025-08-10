package tvshowdeliverystate

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/delivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params delivery.GenerateSearchQueryParams) (string, error)
	SearchTorrent(ctx context.Context, params delivery.SearchTorrentParams) (*delivery.TorrentSearchResult, error)
	GetMagnetLink(ctx context.Context, params delivery.GetMagnetLinkParams) (*delivery.TorrentInfo, error)
	AddTorrentToTorrentClient(ctx context.Context, params delivery.AddTorrentParams) error
	PrepareFileMatches(ctx context.Context, params delivery.PreparingFileMatchesParams) ([]delivery.ContentMatches, error)
	WaitingTorrentDownloadComplete(ctx context.Context, params delivery.WaitingTorrentDownloadCompleteParams) (*delivery.TorrentDownloadStatus, error)
	CreateContentCatalogs(ctx context.Context, params delivery.CreateContentCatalogsParams) (*delivery.TVShowCatalogPath, error)
	StartMergeVideo(ctx context.Context, params delivery.MergeVideoParams) ([]delivery.MergeVideoFile, error)
	GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*delivery.MergeVideoStatus, error)
	SetVideoFileGroup(ctx context.Context, files []string) error
	GetCatalogSize(ctx context.Context, catalogPath string) (uint64, error)
	SetMediaMetaData(ctx context.Context, params delivery.SetMediaMetaDataParams) error
}
