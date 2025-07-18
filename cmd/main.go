package main

import (
	"context"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/storage"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
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

	store := storage.NewStorage()

	lib := tvshowlibrary.NewService(store, themoviedbApi)
	searchResult, err := lib.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
		Query: "Сага о винланде",
	})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("searchResult: %+v\n", searchResult)
}
