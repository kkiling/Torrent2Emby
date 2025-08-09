package common

import (
	"fmt"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
)

type TVShowID struct {
	ID           uint64
	SeasonNumber uint8
}

type ContentID struct {
	MovieID *uint64
	TVShow  *TVShowID
}

func (c *ContentID) Validate() error {
	if c.MovieID == nil && c.TVShow == nil {
		return fmt.Errorf("movieID or tvShow is required: %w", ucerr.InvalidArgument)
	}
	if c.MovieID != nil && c.TVShow != nil {
		return fmt.Errorf("movieID and tvShow cannot be used together: %w", ucerr.InvalidArgument)
	}
	if c.MovieID != nil && *c.MovieID <= 0 {
		return fmt.Errorf("movieID must be positive: %w", ucerr.InvalidArgument)
	}
	if c.TVShow != nil && c.TVShow.ID == uint64(0) {
		return fmt.Errorf("tvShowID must be positive: %w", ucerr.InvalidArgument)
	}
	return nil
}
