package tvshowlibrary

import (
	"context"

	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
)

type TheMovieDb interface {
	SearchTV(ctx context.Context, params themoviedb.SearchQuery) (*themoviedb.TVShowSearchResponse, error)
	GetTV(ctx context.Context, tvID uint64, language themoviedb.Language) (*themoviedb.TVShow, error)
	GetSeasonEpisodes(ctx context.Context, tvID uint64, seasonNumber uint8, language themoviedb.Language) ([]themoviedb.Episode, error)
}

type Storage interface {
	SaveOrUpdateTVShow(ctx context.Context, tvShow *TVShow) error
	GetTVShow(ctx context.Context, tvID uint64) (*TVShow, error)
	GetTVShows(ctx context.Context) ([]TVShowShort, error)
	GetSeasonEpisodes(ctx context.Context, tvID uint64, seasonNumber uint8) ([]Episode, error)
	SaveOrUpdateSeasonEpisode(ctx context.Context, tvID uint64, seasonNumber uint8, episodes []Episode) error
}
