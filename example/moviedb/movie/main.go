package main

import (
	"github.com/kkiling/torrent2emby/internal/apierr"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/themoviedb"
	"time"
)

func main() {
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	api, err := themoviedb.NewApi(
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
		apierr.PrintError(logger, err)
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
		apierr.PrintError(logger, err)
	}
	logger.Info(movieInfo)
}
