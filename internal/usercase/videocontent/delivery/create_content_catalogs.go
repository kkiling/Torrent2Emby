package delivery

import (
	"context"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/samber/lo"

	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

type CreateContentCatalogsParams struct {
	TVShowID common.TVShowID
}

func (s *Service) createTVShowCatalog(ctx context.Context, tvShowID common.TVShowID) (*CatalogsInfo, error) {
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: tvShowID.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return nil, fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	season, find := lo.Find(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == tvShowID.SeasonNumber
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
	seasonName := fmt.Sprintf("S%02d %s", tvShowID.SeasonNumber, season.Name)

	tvShowsPath := filepath.Join(s.config.BasePath, s.config.TVShowMediaSaveTvShowsPath, tvShowName)
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
		if _, err = os.Stat(currentPath); os.IsNotExist(err) {
			// Создаем каталог

			if mkdirErr := syscall.Mkdir(currentPath, 0775); mkdirErr != nil {
				return fmt.Errorf("syscall.Mkdir: %w", mkdirErr)
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

func isEmpty(dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// CreateContentCatalogs формирование каталога куда будет сохранен контент
func (s *Service) CreateContentCatalogs(ctx context.Context, params CreateContentCatalogsParams) (*CatalogsInfo, error) {
	catalog, err := s.createTVShowCatalog(ctx, params.TVShowID)
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}

	if catalog == nil {
		return nil, fmt.Errorf("catalog is nil")
	}

	if createErr := s.createDirectories(s.config.BasePath, catalog.TvShowSeasonPath); createErr != nil {
		return nil, fmt.Errorf("createDirectories: %w", createErr)
	}

	if ok, err := isEmpty(catalog.TvShowSeasonPath); err != nil {
		return nil, fmt.Errorf("isEmpty: %w", err)
	} else if !ok {
		return nil, fmt.Errorf("catalog is not empty: %w", ucerr.AlreadyExists)
	}

	return catalog, nil
}
