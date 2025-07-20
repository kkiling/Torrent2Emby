package main

import (
	"context"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery/deliverystate"

	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	statemachinesqlite "github.com/kkiling/torrent2emby/internal/statemachine/storage/sqlite"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary/storage/sqlite"
)

func main() {
	ctx := context.Background()

	tvShowLibraryStorage, err := sqlite.NewStorage(sqlite.Config{
		DSN: "/home/kiling/projects/torrent2emby/torrent2emby.db",
	}, logger)
	if err != nil {
		logger.Fatal(err)
	}

	stateStorage, err := statemachinesqlite.NewStorage(statemachinesqlite.Config{
		DSN: "/home/kiling/projects/torrent2emby/torrent2emby.db",
	}, logger)

	tvShowLibrary := tvshowlibrary.NewService(tvShowLibraryStorage, themoviedbApi)
	searchResult, err := tvShowLibrary.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
		Query: "Сага о Винланде",
	})
	if err != nil {
		logger.Fatal(err)
	}

	info, err := tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: searchResult.Items[0].ID,
	})
	if err != nil {
		logger.Fatal(err)
	}

	delivery := contentdelivery.NewService(
		contentdelivery.Config{
			BasePath:              "",
			TVShowTorrentSavePath: "",
			TvShowMediaSavePath:   "",
		},
		tvShowLibrary,
		nil,
		nil,
		nil,
		nil)
	deliveryState := deliverystate.NewState(delivery, stateStorage)

	state, err := deliveryState.Create(ctx, deliverystate.CreateOptions{
		MediaID: contentdelivery.MediaID{
			TVShow: &contentdelivery.TVShowID{
				TVShowID:     searchResult.Items[0].ID,
				SeasonNumber: info.Result.Seasons[1].SeasonNumber,
			},
		},
	})
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(state)
}
