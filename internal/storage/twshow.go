package storage

import (
	"context"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

func (s Storage) SaveTVShow(ctx context.Context, tvShow *tvshowlibrary.TVShow) error {
	// TODO implement me
	panic("implement me")
}

func (s Storage) GetTVShow(ctx context.Context, tvID uint64) (*tvshowlibrary.TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetTVShows(ctx context.Context) ([]tvshowlibrary.TVShowShort, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetSeasonEpisodes(ctx context.Context, tvID uint64, seasonNumber int) ([]tvshowlibrary.Episode, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SaveSeasonEpisode(ctx context.Context, tvID uint64, seasonNumber int, episodes []tvshowlibrary.Episode) error {
	//TODO implement me
	panic("implement me")
}
