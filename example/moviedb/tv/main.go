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

	// Поиск сериалов
	response, err := api.SearchTV(themoviedb.SearchQuery{
		Language: themoviedb.LanguageRU,
		Query:    "атака ти",
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

	for _, torrent := range response.Results[:2] {
		logger.Infof("#%d - %s - %s (%f%d)(%f)",
			torrent.ID,
			torrent.Name,
			torrent.FirstAirDate.Format(time.DateOnly),
			torrent.VoteAverage,
			torrent.VoteCount,
			torrent.Popularity)
	}

	// информация о сериале
	tvId := response.Results[0].ID
	tvOInfo, err := api.GetTV(tvId, themoviedb.LanguageRU)
	if err != nil {
		apierr.PrintError(logger, err)
	}
	logger.Info(tvOInfo)

	episodes, err := api.GetSeasonEpisodes(tvId, 1, themoviedb.LanguageRU)
	if err != nil {
		apierr.PrintError(logger, err)
	}
	for _, episode := range episodes {
		logger.Infof("%s (%s)", episode.Name, episode.AirDate)
	}
}
