package tvshowdeliverystate

import (
	"context"

	"github.com/google/uuid"

	delivery2 "github.com/kkiling/torrent2emby/internal/usercase/videocontent/delivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params delivery2.GenerateSearchQueryParams) (string, error)
	SearchTorrent(ctx context.Context, params delivery2.SearchTorrentParams) (*delivery2.TorrentSearchResult, error)
	GetMagnetLink(ctx context.Context, params delivery2.GetMagnetLinkParams) (*delivery2.MagnetInfo, error)
	AddTorrentToTorrentClient(ctx context.Context, params delivery2.AddTorrentParams) error
	PrepareFileMatches(ctx context.Context, params delivery2.PreparingFileMatchesParams) ([]delivery2.ContentMatches, error)
	WaitingTorrentDownloadComplete(ctx context.Context, params delivery2.WaitingTorrentDownloadCompleteParams) (*delivery2.TorrentDownloadStatus, error)
	CreateContentCatalogs(ctx context.Context, params delivery2.CreateContentCatalogsParams) (*delivery2.CatalogsInfo, error)
	StartMergeVideo(ctx context.Context, params delivery2.MergeVideoParams) ([]delivery2.MergeVideoFile, error)
	GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*delivery2.MergeVideoStatus, error)
	SetVideoFileGroup(ctx context.Context, files []string) error
	SetMediaMetaData(ctx context.Context, params delivery2.SetMediaMetaDataParams) error
}
