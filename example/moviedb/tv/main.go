package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
	themoviedb2 "github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	api, err := themoviedb2.NewApi(
		logger,
		cfg.MovieDb.ApiKey,
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Поиск сериалов
	response, err := api.SearchTV(ctx, themoviedb2.SearchQuery{
		Language: themoviedb2.LanguageRU,
		Query:    "Сага о винланде",
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

	for _, tv := range response.Results[:1] {
		logger.Infof("#%d - %s - %s (%f%d)(%f)",
			tv.ID,
			tv.Name,
			tv.FirstAirDate.Format(time.DateOnly),
			tv.VoteAverage,
			tv.VoteCount,
			tv.Popularity)
	}

	// информация о сериале
	tvId := response.Results[0].ID
	tvOInfo, err := api.GetTV(ctx, tvId, themoviedb2.LanguageEN)
	if err != nil {
		apierr.PrintError(logger, err)
	}
	logger.Info(tvOInfo.Seasons[2].Name)

	episodes, err := api.GetSeasonEpisodes(ctx, tvId, 2, themoviedb2.LanguageEN)
	if err != nil {
		apierr.PrintError(logger, err)
	}
	for _, episode := range episodes {
		fmt.Printf("#%d %s\n", episode.EpisodeNumber, episode.Name)
	}
}
