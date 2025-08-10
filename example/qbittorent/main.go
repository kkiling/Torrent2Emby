package main2

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/kkiling/goplatform/config"
	"github.com/kkiling/goplatform/log"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	appconfig "github.com/kkiling/torrent2emby/internal/config"
)

func main() {
	logger := log.NewLogger(log.DebugLevel)

	var args config.Arguments
	if _, err := flags.Parse(&args); err != nil {
		logger.Fatal(err)
	}
	cfgProvider, err := config.NewProvider(args)
	if err != nil {
		logger.Fatal(err)
	}
	cfg, err := appconfig.NewEnvConfig(cfgProvider)
	if err != nil {
		logger.Fatal(err)
	}

	// Создаем сервис для работы с рутрекером
	api, err := qbittorrent.NewApi(
		logger,
		cfg.QBittorrent.ApiUrl,
		cfg.QBittorrent.Username,
		cfg.QBittorrent.Password,
		cfg.QBittorrent.CookieDir,
	)
	if err != nil {
		logger.Fatal(err)
	}

	//// Делаем запрос на поиск раздач по названию
	//err = api.AddTorrent(qbittorrent.TorrentAddOptions{
	//	// hash 91010E2E27CD888EA7013DF2AD152F01DB0437DF
	//	Magnet:   "magnet:?xt=urn:btih:91010E2E27CD888EA7013DF2AD152F01DB0437DF&tr=http%3A%2F%2Fbt.t-ru.org%2Fann%3Fmagnet",
	//	SavePath: "/downloads/tvshows",
	//	Category: "anime",
	//	Tags:     []string{"anime", "tvshows"},
	//	Paused:   true,
	//})
	//if err != nil {
	//	apierr.PrintError(logger, err)
	//}

	//err = api.DeleteTorrent("91010E2E27CD888EA7013DF2AD152F01DB0437DF", true)
	//if err != nil {
	//	apierr.PrintError(logger, err)
	//}

	//err = api.ResumeTorrent("91010E2E27CD888EA7013DF2AD152F01DB0437DF")
	//if err != nil {
	//	apierr.PrintError(logger, err)
	//}

	//err = api.PauseTorrent("91010E2E27CD888EA7013DF2AD152F01DB0437DF")
	//if err != nil {
	//	apierr.PrintError(logger, err)
	//}
	//
	//info, err := api.GetTorrentInfo("91010E2E27CD888EA7013DF2AD152F01DB0437DF")
	//if err != nil {
	//	apierr.PrintError(logger, err)
	//}
	//logger.Infof("Torrent info")
	//logger.Infof("Name: %s", info.Name)
	//logger.Infof("State: %s", info.State)

	files, err := api.GetTorrentFiles("f1770424dc4ccf56a466c61ff92044c868b70645")
	if err != nil {
		apierr.PrintError(logger, err)
	}

	logger.Infof("Torrent files")
	for _, file := range files {
		// logger.Infof("%s (%d) - %f", prepare.Name, prepare.Size, prepare.Progress)
		fmt.Println(file.Name)
	}

	logger.Info("Done")
}
