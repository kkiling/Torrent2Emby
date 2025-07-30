package deliverystate

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params contentdelivery.GenerateSearchQueryParams) (string, error)
	SearchTorrent(ctx context.Context, params contentdelivery.SearchTorrentParams) (*contentdelivery.TorrentSearchResult, error)
	GetMagnetLink(ctx context.Context, params contentdelivery.GetMagnetLinkParams) (*contentdelivery.MagnetInfo, error)
	AddTorrentToTorrentClient(ctx context.Context, params contentdelivery.AddTorrentParams) error
	PrepareFileMatches(ctx context.Context, params contentdelivery.PreparingFileMatchesParams) ([]contentdelivery.ContentMatches, error)
	WaitingTorrentDownloadComplete(ctx context.Context, params contentdelivery.WaitingTorrentDownloadCompleteParams) (*contentdelivery.TorrentDownloadStatus, error)
	CreateContentCatalogs(ctx context.Context, params contentdelivery.CreateContentCatalogsParams) (*contentdelivery.CatalogsInfo, error)
	StartMergeVideo(ctx context.Context, params contentdelivery.MergeVideoParams) ([]contentdelivery.MergeVideoFile, error)
	GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*contentdelivery.MergeVideoStatus, error)
	SetVideoFileGroup(ctx context.Context, files []string) error
	SetMediaMetaData(ctx context.Context, params contentdelivery.SetMediaMetaDataParams) error
}
