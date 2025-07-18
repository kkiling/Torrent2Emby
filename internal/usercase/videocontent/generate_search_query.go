package videocontent

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

type GenerateSearchQueryParams struct {
	MediaID MediaID
}

func (s *Service) getTVShowQuery(ctx context.Context, tvShowID uint64, seasonNumber int) (string, error) {
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: tvShowID,
	})
	if err != nil {
		return "", fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return "", fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	season, find := lo.Find(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == seasonNumber
	})
	if !find {
		return "", fmt.Errorf("season not found: %w", ucerr.NotFound)
	}
	// Формируем поисковый запрос на основе инфы  о сезоне сериала
	searchQuery := fmt.Sprintf("%s сезон %d", tvShowInfo.Result.Name, season.SeasonNumber)

	return searchQuery, nil
}

// generateSearchQuery Формируем поисковый запрос к торент трекеру на основе данных сезона сериала / фильма
func (s *Service) generateSearchQuery(ctx context.Context, params GenerateSearchQueryParams) (string, error) {
	searchQuery := ""
	if params.MediaID.TVShow != nil {
		var err error
		searchQuery, err = s.getTVShowQuery(ctx, params.MediaID.TVShow.TVShowID, params.MediaID.TVShow.SeasonNumber)
		if err != nil {
			return "", fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
		}
	}
	if params.MediaID.MovieID != nil {
		return "", fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	return searchQuery, nil
}
