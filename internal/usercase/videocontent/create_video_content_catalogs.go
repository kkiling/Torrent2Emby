package videocontent

import (
	"context"
	"fmt"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

type CreateVideoContentCatalogsParams struct {
	MediaID MediaID
}

func (s *Service) createTVShowCatalog(ctx context.Context, tvShowID uint64, seasonNumber int) (string, error) {
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
	// Формируем каталог
	tvShowName := fmt.Sprintf("%s (%d)", tvShowInfo.Result.Name, tvShowInfo.Result.FirstAirDate.Year())
	seasonName := fmt.Sprintf("#%d %s", seasonNumber, season.Name)
	result := filepath.Join(s.config.BasePath, tvShowName, seasonName)

	return result, nil
}

// createVideoContentCatalogs формирование каталога куда будет сохранен контент
func (s *Service) createVideoContentCatalogs(ctx context.Context, params CreateVideoContentCatalogsParams) (CatalogsInfo, error) {
	var catalog = ""
	if params.MediaID.TVShow != nil {
		var err error
		catalog, err = s.createTVShowCatalog(ctx, params.MediaID.TVShow.TVShowID, params.MediaID.TVShow.SeasonNumber)
		if err != nil {
			return CatalogsInfo{}, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
		}
	}
	if params.MediaID.MovieID != nil {
		return CatalogsInfo{}, fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	// Создание необходимых каталогов
	err := os.MkdirAll(catalog, os.ModeDir)
	if err != nil {
		return CatalogsInfo{}, fmt.Errorf("os.MkdirAll: %w", err)
	}

	return CatalogsInfo{CatalogPath: catalog}, nil
}
