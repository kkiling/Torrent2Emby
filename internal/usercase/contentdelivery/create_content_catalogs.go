package contentdelivery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/samber/lo"

	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

type CreateContentCatalogsParams struct {
	MediaID MediaID
}

func (s *Service) createTVShowCatalog(ctx context.Context, tvShowID uint64, seasonNumber int) (*CatalogsInfo, error) {
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: tvShowID,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return nil, fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	season, find := lo.Find(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == seasonNumber
	})
	if !find {
		return nil, fmt.Errorf("season not found: %w", ucerr.NotFound)
	}
	// Формируем каталог
	// Название сезона
	/*
		Series Name/
		  Season 01/
		    S01E01 - Episode Name.mp4
	*/
	tvShowName := fmt.Sprintf("%s (%d)", tvShowInfo.Result.Name, tvShowInfo.Result.FirstAirDate.Year())
	seasonName := fmt.Sprintf("S%02d %s", seasonNumber, season.Name)

	tvShowsPath := filepath.Join(s.config.BasePath, s.config.TvShowMediaSaveTvShowsPath, tvShowName)
	return &CatalogsInfo{
		TvShowPath:       tvShowsPath,
		TvShowSeasonPath: filepath.Join(tvShowsPath, seasonName),
	}, nil
}

func (s *Service) createDirectories(base, catalog string) error {
	// Проверяем, что catalog действительно является подкаталогом base
	relPath, err := filepath.Rel(base, catalog)
	if err != nil {
		return fmt.Errorf("catalog is not a subdirectory of base: %v", err)
	}

	// Разбиваем относительный путь на компоненты
	parts := strings.Split(relPath, string(filepath.Separator))

	// Постепенно создаём каталоги
	currentPath := base
	for _, part := range parts {
		currentPath = filepath.Join(currentPath, part)
		// Проверяем существование каталога
		if _, err := os.Stat(currentPath); os.IsNotExist(err) {
			// Создаем каталог
			err := syscall.Mkdir(currentPath, 0775)
			if err != nil {
				return fmt.Errorf("syscall.Mkdir: %w", err)
			}
			// Меняем группу пользователей
			if s.config.UserGroup != "" {
				if err = setGroup(currentPath, s.config.UserGroup); err != nil {
					return fmt.Errorf("syscall.Chown: %w", err)
				}
			}
		} else if err != nil {
			return fmt.Errorf("error checking directory %s: %v", currentPath, err)
		}
	}

	return nil
}

// CreateContentCatalogs формирование каталога куда будет сохранен контент
func (s *Service) CreateContentCatalogs(ctx context.Context, params CreateContentCatalogsParams) (*CatalogsInfo, error) {
	var catalog *CatalogsInfo
	if params.MediaID.TVShow != nil {
		var err error
		catalog, err = s.createTVShowCatalog(ctx, params.MediaID.TVShow.TVShowID, params.MediaID.TVShow.SeasonNumber)
		if err != nil {
			return nil, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
		}
	}
	if params.MediaID.MovieID != nil {
		return nil, fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	if catalog == nil {
		return nil, fmt.Errorf("catalog is nil")
	}

	err := s.createDirectories(s.config.BasePath, catalog.TvShowSeasonPath)
	if err != nil {
		return nil, fmt.Errorf("createDirectories: %w", err)
	}

	return catalog, nil
}
