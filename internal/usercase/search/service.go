package search

import (
	"context"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/themoviedb"
)

const (
	language       = themoviedb.LanguageRU
	perPageDefault = 20
)

type TheMovieDb interface {
	SearchMovie(params themoviedb.SearchQuery) (*themoviedb.MovieSearchResponse, error)
	SearchTV(params themoviedb.SearchQuery) (*themoviedb.TVShowSearchResponse, error)
	GetTV(tvID uint64, language themoviedb.Language) (*themoviedb.TVShow, error)
	GetSeasonEpisodes(tvID uint64, seasonNumber int, language themoviedb.Language) ([]themoviedb.Episode, error)
}

type Service struct {
	theMovieDb TheMovieDb
}

func NewSearch(
	theMovieDb TheMovieDb,
) *Service {
	return &Service{
		theMovieDb: theMovieDb,
	}
}

// SearchMovie поиск фильмов по названию
func (s *Service) SearchMovie(_ context.Context, params TVShowSearchParams) (*MovieSearchResult, error) {
	response, err := s.theMovieDb.SearchMovie(themoviedb.SearchQuery{
		Language: language,
		Query:    params.Query,
		Page:     1,
		PerPage:  perPageDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.SearchMovie: %w", err)
	}

	return mapMovieSearchResult(response)
}

// SearchTVShow поиск сериалов по названию
func (s *Service) SearchTVShow(_ context.Context, params TVShowSearchParams) (*TVShowSearchResult, error) {
	response, err := s.theMovieDb.SearchTV(themoviedb.SearchQuery{
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
	response, err := s.theMovieDb.GetTV(params.TVShowID, language)
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.GetTV: %w", err)
	}

	return mapGetTVShowResult(response)
}

// GetSeasonEpisodes получение информации о сериях сезона
func (s *Service) GetSeasonEpisodes(_ context.Context, params GetSeasonEpisodesParams) (*GetSeasonEpisodesResult, error) {
	response, err := s.theMovieDb.GetSeasonEpisodes(params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.GetSeasonEpisodes: %w", err)
	}

	return mapGetSeasonEpisodesResult(response)
}

// AddTVShowToLibrary добавление сериала в библиотеку
func (s *Service) AddTVShowToLibrary(_ context.Context, params AddTVShowToLibraryParams) error {
	panic("implement me")
}

// GetTVShowsFromLibrary получение списка сериалов из библиотеки
func (s *Service) GetTVShowsFromLibrary(_ context.Context, params GetTVShowsFromLibraryParams) (*GetTVShowsFromLibraryResult, error) {
	panic("implement me")
}
