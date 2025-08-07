package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/kkiling/goplatform/config"
	"github.com/kkiling/goplatform/log"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
	"github.com/kkiling/torrent2emby/internal/adapter/rutracker"
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
	rutrackerApi, err := rutracker.NewApi(
		logger,
		cfg.Rutracker.Username,
		cfg.Rutracker.Password,
		cfg.Rutracker.CookieDir,
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Делаем запрос на поиск раздач по названию
	response, err := rutrackerApi.SearchTorrents("клинок рассекающий демонов ")
	if err != nil {
		apierr.PrintError(logger, err)
	}

	// Может быть так что раздачи не найдены
	if len(response.Results) == 0 {
		logger.Info("No results")
		return
	}

	for _, torrent := range response.Results[:3] {
		logger.Infof("%s (%s)", torrent.Title, torrent.Size)
	}

	// По первой найденной раздачи пытаемся получить magnet ссылку
	magnet, err := rutrackerApi.GetMagnetLink(response.Results[0].Href)
	if err != nil {
		apierr.PrintError(logger, err)
	}

	logger.Infof("Magnet: %s", magnet.Magnet)
	logger.Infof("Hash: %s", magnet.Hash)
	logger.Info("Done")
}
