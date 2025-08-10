package delivery

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
)

type GenerateSearchQueryParams struct {
	TVShowID common.TVShowID
}

func (s *Service) getTVShowQuery(ctx context.Context, tvShowID common.TVShowID) (string, error) {
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: tvShowID.ID,
	})
	if err != nil {
		return "", fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return "", fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	season, find := lo.Find(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == tvShowID.SeasonNumber
	})
	if !find {
		return "", fmt.Errorf("season not found: %w", ucerr.NotFound)
	}
	// Формируем поисковый запрос на основе инфы о сезоне сериала
	searchQuery := fmt.Sprintf("%s %d", tvShowInfo.Result.Name, season.AirDate.Year())

	return searchQuery, nil
}

// GenerateSearchQuery формируем поисковый запрос к торент трекеру на основе данных сезона сериала / фильма
func (s *Service) GenerateSearchQuery(ctx context.Context, params GenerateSearchQueryParams) (string, error) {
	searchQuery, err := s.getTVShowQuery(ctx, params.TVShowID)
	if err != nil {
		return "", fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}

	return searchQuery, nil
}
