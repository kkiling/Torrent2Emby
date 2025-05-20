package main

import (
	"errors"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/themoviedb"
	"time"
)

func printError(logger log.Logger, err error) {
	switch {
	case errors.Is(err, themoviedb.NotAuthorizedErr):
		logger.Fatal("Клиент не авторизован")
	case errors.Is(err, themoviedb.AuthenticationFailedErr):
		logger.Fatal("Ошибка при попытке залогиниться")
	case errors.Is(err, themoviedb.ServiceUnavailableErr): // предполагаемое название ошибки
		logger.Fatal("Сервис не доступен")
	default:
		logger.Fatal(err)
	}
}

func main() {
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	api, err := themoviedb.NewAPI(
		logger,
		cfg.MovieDb.ApiKey,
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Делаем запрос на поиск раздач по названию
	response, err := api.SearchMovie(themoviedb.SearchQuery{
		Language: themoviedb.LanguageRU,
		Query:    "бойцовский клуб",
		Page:     1,
		PerPage:  10,
	})
	if err != nil {
		printError(logger, err)
	}

	if len(response.Results) == 0 {
		logger.Info("No results")
		return
	}

	for _, torrent := range response.Results[:3] {
		logger.Infof("#%d - %s - %s (%f%d)(%f)",
			torrent.ID,
			torrent.Title,
			torrent.ReleaseDate.Format(time.DateOnly),
			torrent.VoteAverage,
			torrent.VoteCount,
			torrent.Popularity)
	}

	movieId := response.Results[0].ID
	movieInfo, err := api.GetMovie(movieId, themoviedb.LanguageRU)
	if err != nil {
		printError(logger, err)
	}
	logger.Info(movieInfo)
}
