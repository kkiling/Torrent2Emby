package main

import (
	"errors"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/rutracker"
)

func printError(logger log.Logger, err error) {
	switch {
	case errors.Is(err, rutracker.NotAuthorizedErr):
		logger.Fatal("Клиент не авторизован")
	case errors.Is(err, rutracker.AuthenticationFailedErr):
		logger.Fatal("Ошибка при попытке залогиниться")
	case errors.Is(err, rutracker.ServiceUnavailableErr): // предполагаемое название ошибки
		logger.Fatal("Сервис не доступен")
	default:
		logger.Fatal(err)
	}
}

func main() {
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Создаем сервис для работы с рутрекером
	rutrackerApi, err := rutracker.NewAPI(
		logger,
		cfg.Rutracker.Username,
		cfg.Rutracker.Password,
		cfg.Rutracker.CookieDir,
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Делаем запрос на поиск раздач по названию
	response, err := rutrackerApi.SearchTorrents("бойцовский клуб")
	if err != nil {
		printError(logger, err)
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
		printError(logger, err)
	}

	logger.Infof("Magnet: %s", magnet.Magnet)
	logger.Infof("Hash: %s", magnet.Hash)
	logger.Info("Done")
}
