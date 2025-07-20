package videodelivery

import (
	"context"

	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	"github.com/kkiling/torrent2emby/internal/adapter/prepare/tvshow"
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
}

type PrepareTVShow interface {
	PrepareTvShowSeason(params *tvshow.PrepareTvShowPrams) (*tvshow.PrepareTVShowSeason, error)
}

type MkvMerge interface {
	Merge(params mkvmerge.MergeParams) error
	GetMediaInfo(filePath string) (*mkvmerge.MediaInfo, error)
}
