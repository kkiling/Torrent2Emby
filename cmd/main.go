package main

import (
	"context"
	"fmt"

	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary/storage/sqlite"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	themoviedbApi, err := themoviedb.NewApi(
		logger,
		cfg.MovieDb.ApiKey,
	)
	if err != nil {
		logger.Fatal(err)
	}

	tvShowLibraryStorage, err := sqlite.NewStorage(sqlite.Config{
		DSN: "/home/kiling/projects/torrent2emby/torrent2emby.db",
	}, logger)
	if err != nil {
		logger.Fatal(err)
	}

	lib := tvshowlibrary.NewService(tvShowLibraryStorage, themoviedbApi)
	searchResult, err := lib.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
		Query: "Сага о Винланде",
	})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("searchResult: %+v\n", searchResult)

	info, err := lib.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: searchResult.Items[0].ID,
	})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("info: %+v\n", info)

	episodes, err := lib.GetSeasonEpisodes(ctx, tvshowlibrary.GetSeasonEpisodesParams{
		TVShowID:     searchResult.Items[0].ID,
		SeasonNumber: info.Result.Seasons[1].SeasonNumber,
	})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("seasonn: %+v\n", episodes)

}
