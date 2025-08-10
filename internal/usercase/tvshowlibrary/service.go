package tvshowlibrary

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkiling/goplatform/storagebase"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
)

const (
	language       = themoviedb.LanguageRU
	perPageDefault = 20
)

type Service struct {
	theMovieDb TheMovieDb
	storage    Storage
}

func NewService(
	storage Storage,
	theMovieDb TheMovieDb,
) *Service {
	return &Service{
		storage:    storage,
		theMovieDb: theMovieDb,
	}
}

// SearchTVShow поиск сериалов по названию
func (s *Service) SearchTVShow(ctx context.Context, params TVShowSearchParams) (*TVShowSearchResult, error) {
	response, err := s.theMovieDb.SearchTV(ctx, themoviedb.SearchQuery{
		Language: language,
		Query:    params.Query,
		Page:     1,
		PerPage:  perPageDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.SearchTV: %w", err)
	}

	return &TVShowSearchResult{
		Items: mapTVShowShorts(response.Results),
	}, nil
}

// GetTVShowInfo получение подробной информации о сериале
func (s *Service) GetTVShowInfo(ctx context.Context, params GetTVShowParams) (*GetTVShowResult, error) {
	// Сначала тянем информацию о сериале из библиотеки
	if tvShow, err := s.storage.GetTVShow(ctx, params.TVShowID); err != nil {
		switch {
		case errors.Is(err, storagebase.ErrNotFound):
			// Не найдено, идем дальше
		default:
			return nil, fmt.Errorf("storage.GetStateByIdempotencyKey: %w", err)
		}
	} else {
		return &GetTVShowResult{
			Result: tvShow,
		}, err
	}

	response, err := s.theMovieDb.GetTV(ctx, params.TVShowID, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return nil, ucerr.NotFound
		}
		return nil, fmt.Errorf("theMovieDb.GetTV: %w", err)
	}

	tvShow := mapTVShow(response)

	// Получение информации о сериале, автоматически добавляет его в библиотеку
	if err = s.storage.SaveOrUpdateTVShow(ctx, tvShow); err != nil {
		return nil, fmt.Errorf("s.storage.SaveTVShow: %w", err)
	}

	return &GetTVShowResult{
		Result: tvShow,
	}, err
}

// GetSeasonEpisodes получение информации о сериях сезона
func (s *Service) GetSeasonEpisodes(ctx context.Context, params GetSeasonEpisodesParams) (*GetSeasonEpisodesResult, error) {
	// Сначала проверяем существует ли сам сериал
	if _, err := s.GetTVShowInfo(ctx, GetTVShowParams{TVShowID: params.TVShowID}); err != nil {
		return nil, fmt.Errorf("s.GetTVShowInfo: %w", err)
	}

	// сначала тянем информацию о эпизодах из библиотеки
	if episodes, err := s.storage.GetSeasonEpisodes(ctx, params.TVShowID, params.SeasonNumber); err != nil {
		switch {
		case errors.Is(err, storagebase.ErrNotFound):
			// Не найдено, идем дальше
		default:
			return nil, fmt.Errorf("storage.GetTVShow: %w", err)
		}
	} else if len(episodes) > 0 {
		return &GetSeasonEpisodesResult{
			Items: episodes,
		}, err
	}

	response, err := s.theMovieDb.GetSeasonEpisodes(ctx, params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return nil, ucerr.NotFound
		}
		return nil, fmt.Errorf("theMovieDb.GetSeasonEpisodes: %w", err)
	}

	episodes := mapEpisodes(response)

	// Получение информации о эпизодах сериала, автоматически добавляет их в библиотеку
	if err = s.storage.SaveOrUpdateSeasonEpisode(ctx, params.TVShowID, params.SeasonNumber, episodes); err != nil {
		return nil, fmt.Errorf("s.storage.SaveSeasonEpisode: %w", err)
	}

	return &GetSeasonEpisodesResult{
		Items: mapEpisodes(response),
	}, nil
}

// GetTVShowsFromLibrary получение списка сериалов из библиотеки
func (s *Service) GetTVShowsFromLibrary(ctx context.Context, _ GetTVShowsFromLibraryParams) (*GetTVShowsFromLibraryResult, error) {
	tvShows, err := s.storage.GetTVShows(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.GetTVShow: %w", err)
	}
	return &GetTVShowsFromLibraryResult{
		Items: tvShows,
	}, nil
}
