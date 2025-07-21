package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/container"
	"github.com/kkiling/torrent2emby/internal/statemachine"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery/deliverystate"
	"log"

	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

func main() {
	ctx := context.Background()

	cn, err := container.NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	tvShowLibrary := cn.GetTvShowLibrary()
	deliveryStateMachine := cn.DeliveryStateMachine()

	searchResult, err := tvShowLibrary.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
		Query: "Сага о Винланде",
	})
	if err != nil {
		log.Fatal(err)
	}

	info, err := tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: searchResult.Items[0].ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	state, err := deliveryStateMachine.Create(ctx, deliverystate.CreateOptions{
		MediaID: contentdelivery.MediaID{
			TVShow: &contentdelivery.TVShowID{
				TVShowID:     searchResult.Items[0].ID,
				SeasonNumber: info.Result.Seasons[2].SeasonNumber,
			},
		},
	})
	if err != nil {
		if !errors.Is(err, statemachine.ErrAlreadyExists) {
			log.Fatal(err)
		}
	}

	fmt.Println(state)

	//newState, eerr, err := deliveryStateMachine.Complete(ctx, state.ID, deliverystate.ChoseTorrentOptions{
	//	//NewSearchQuery: lo.ToPtr("Сага о Винланде 2023"),
	//	Href: lo.ToPtr("https://rutracker.org/forum/viewtopic.php?t=6313846"),
	//})
	//newState, eerr, err := deliveryStateMachine.Complete(ctx, state.ID, deliverystate.ChoseFileMatchesOptions{
	//	Approve: true,
	//})
	newState, eerr, err := deliveryStateMachine.Complete(ctx, state.ID)
	if err != nil {
		log.Fatal(err)
	}
	if eerr != nil {
		log.Fatal(eerr)
	}
	fmt.Println(newState)

}
