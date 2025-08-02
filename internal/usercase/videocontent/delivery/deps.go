package delivery

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/torrent2emby/internal/adapter/emby"
	tvshow2 "github.com/kkiling/torrent2emby/internal/adapter/matchtvshow"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/adapter/rutracker"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

type TVShowLibrary interface {
	GetTVShowInfo(ctx context.Context, params tvshowlibrary.GetTVShowParams) (*tvshowlibrary.GetTVShowResult, error)
	GetSeasonEpisodes(ctx context.Context, params tvshowlibrary.GetSeasonEpisodesParams) (*tvshowlibrary.GetSeasonEpisodesResult, error)
}

type TorrentSite interface {
	SearchTorrents(query string) (*rutracker.TorrentResponse, error)
	GetMagnetLink(torrentUrl string) (*rutracker.MagnetInfo, error)
}

type TorrentClient interface {
	AddTorrent(opts qbittorrent.TorrentAddOptions) error
	GetTorrentInfo(hash string) (*qbittorrent.TorrentInfo, error)
	GetTorrentFiles(hash string) ([]qbittorrent.TorrentFile, error)
	ResumeTorrent(hash string) error
}

type PrepareTVShow interface {
	PrepareTvShowSeason(params *tvshow2.PrepareTvShowPrams) (*tvshow2.PrepareTVShowSeason, error)
}

type MkvMergePipeline interface {
	AddToMerge(ctx context.Context, idempotencyKey string, params mkvmerge.MergeParams) (*mkvmerge.MergeResult, error)
	GetMergeResult(ctx context.Context, id uuid.UUID) (*mkvmerge.MergeResult, error)
}

type EmbyApi interface {
	RemoteSearchApply(embyID, theMovieDBID uint64) error
	GetCatalogInfo(path string) (*emby.CatalogInfo, error)
}
