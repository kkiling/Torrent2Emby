package tvshowlibrary

import (
	"context"
	"fmt"
	themoviedb2 "github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
)

const (
	language       = themoviedb2.LanguageRU
	perPageDefault = 20
)

type TheMovieDb interface {
	SearchMovie(params themoviedb2.SearchQuery) (*themoviedb2.MovieSearchResponse, error)
	SearchTV(params themoviedb2.SearchQuery) (*themoviedb2.TVShowSearchResponse, error)
	GetTV(tvID uint64, language themoviedb2.Language) (*themoviedb2.TVShow, error)
	GetSeasonEpisodes(tvID uint64, seasonNumber int, language themoviedb2.Language) ([]themoviedb2.Episode, error)
}

type Service struct {
	theMovieDb TheMovieDb
}

func NewService(
	theMovieDb TheMovieDb,
) *Service {
	return &Service{
		theMovieDb: theMovieDb,
	}
}

// SearchTVShow поиск сериалов по названию
func (s *Service) SearchTVShow(_ context.Context, params TVShowSearchParams) (*TVShowSearchResult, error) {
	response, err := s.theMovieDb.SearchTV(themoviedb2.SearchQuery{
		Language: language,
		Query:    params.Query,
		Page:     1,
		PerPage:  perPageDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.SearchTV: %w", err)
	}

	return mapTvShowSearchResult(response)
}

// GetTVShowInfo получение подробной информации о сериале
func (s *Service) GetTVShowInfo(_ context.Context, params GetTVShowParams) (*GetTVShowResult, error) {
	// TODO: сначала тянем информацию о сериале из библиотеки

	response, err := s.theMovieDb.GetTV(params.TVShowID, language)
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.GetTV: %w", err)
	}
	// TODO: Получение информации о сериале, автоматически добавляет его в библиотеку
	return mapGetTVShowResult(response)
}

// GetSeasonEpisodes получение информации о сериях сезона
func (s *Service) GetSeasonEpisodes(_ context.Context, params GetSeasonEpisodesParams) (*GetSeasonEpisodesResult, error) {
	// TODO: сначала тянем информацию о эпизодах из библиотеки

	response, err := s.theMovieDb.GetSeasonEpisodes(params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.GetSeasonEpisodes: %w", err)
	}
	// TODO: получение информации о эпизодах сериала, автоматически добавляет их в библиотеку
	return mapGetSeasonEpisodesResult(response)
}

// GetTVShowsFromLibrary получение списка сериалов из библиотеки
func (s *Service) GetTVShowsFromLibrary(_ context.Context, params GetTVShowsFromLibraryParams) (*GetTVShowsFromLibraryResult, error) {
	// TODO: По фильтрам ищем сериалы добавленные в библиотеку
	panic("implement me")
}
