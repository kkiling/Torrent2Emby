package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	platformserver "github.com/kkiling/goplatform/server"

	"github.com/kkiling/torrent2emby/internal/container"
	"github.com/kkiling/torrent2emby/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cn, err := container.NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = cn.MkvMergePipeline().StartMergePipeline(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	logger := cn.GetLogger()

	srv := server.NewTorrent2EmbyServer(
		logger,
		platformserver.Config{
			Host:                    "localhost",
			GrpcPort:                8181,
			HttpPort:                8080,
			MaxSendMessageLength:    2147483647,
			MaxReceiveMessageLength: 63554432,
			ShutdownTimeout:         3,
		},
		cn.GetTvShowLibrary(),
	)
	go func() {
		err = srv.Start(ctx)
		if err != nil {
			log.Fatalf("fail start app: %v", err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Infof("--- shutdown application ---")
		cancel()
	}()

	<-ctx.Done()
	logger.Infof("--- stopped application ---")
	srv.Stop()
	logger.Infof("--- stop application ---")

}

//tvShowLibrary := cn.GetTvShowLibrary()
//tvShowDeliveryStateMachine := cn.TVShowDeliveryStateMachine()
//
//searchResult, err := tvShowLibrary.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
//	Query: "Сага о Винланде",
//})
//if err != nil {
//	log.Fatal(err)
//}
//
//info, err := tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
//	TVShowID: searchResult.Items[0].ID,
//})
//if err != nil {
//	log.Fatal(err)
//}
//
//state, err := tvShowDeliveryStateMachine.Create(ctx, tvshowdeliverystate.CreateOptions{
//	TVShowID: videocontent.TVShowID{
//		ID:           searchResult.Items[0].ID,
//		SeasonNumber: info.Result.Seasons[2].SeasonNumber,
//	},
//})
//if err != nil {
//	if !errors.Is(err, statemachine.ErrAlreadyExists) {
//		log.Fatal(err)
//	}

//newState, eerr, err := deliveryStateMachine.Complete(ctx, state.ID, deliverystate.ChoseTorrentOptions{
//	//NewSearchQuery: lo.ToPtr("Сага о Винланде 2023"),
//	Href: lo.ToPtr("https://rutracker.org/forum/viewtopic.php?t=6313846"),
//})
//newState, eerr, err := deliveryStateMachine.Complete(ctx, state.ID, deliverystate.ChoseFileMatchesOptions{
//	Approve: true,
//})

//newState, eerr, err := tvShowDeliveryStateMachine.Complete(ctx, state.ID)
//
//if err != nil {
//	log.Fatal(err)
//}
//if eerr != nil {
//	log.Fatal(eerr)
//}
//
//fmt.Println(newState.Status)
//
//time.Sleep(1 * time.Hour)
