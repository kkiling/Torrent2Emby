package main

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/apierr"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	prepare "github.com/kkiling/torrent2emby/internal/prepare/tvshow"
	"github.com/kkiling/torrent2emby/internal/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/themoviedb"
	"github.com/samber/lo"
	"path/filepath"
)

const (
	torrentHash = "f1770424dc4ccf56a466c61ff92044c868b70645"
	tvShowID    = 88803 // Сага о винланде
	season      = 2
)

func mapToPrepareTvShowPrams(
	savePath, contentPath string,
	episodes []themoviedb.Episode,
	torrentFiles []qbittorrent.TorrentFile,
) (*prepare.PrepareTvShowPrams, error) {

	// Вычисляем относительный путь от savePath до currentPath
	// SavePath: /downloads
	// ContentPath /downloads/Vinland Sag
	// relPath Vinland Sag
	relPath, err := filepath.Rel(savePath, contentPath)
	if err != nil {
		return nil, fmt.Errorf("filepath.Rel: %w", err)
	}

	// Получаем относительный путь файла
	var prepareTorrentFiles []prepare.TorrentFile
	for _, file := range torrentFiles {
		relFile, err := filepath.Rel(relPath, file.Name)
		if err != nil {
			return nil, fmt.Errorf("filepath.Rel: %w", err)
		}
		prepareTorrentFiles = append(prepareTorrentFiles, prepare.TorrentFile{
			RelativePath: relFile,
			Size:         file.Size,
		})
	}

	return &prepare.PrepareTvShowPrams{
		Episodes: lo.Map(episodes, func(episode themoviedb.Episode, _ int) prepare.Episode {
			return prepare.Episode{
				EpisodeNumber: episode.EpisodeNumber,
				Name:          episode.Name,
			}
		}),
		TorrentFiles: prepareTorrentFiles,
	}, nil
}

func main() {
	logger := log.NewLogger(log.DebugLevel)

	cfg, err := config.NewEnvConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Создаем сервис для работы с рутрекером
	torrentApi, err := qbittorrent.NewApi(
		logger,
		cfg.QBittorrent.ApiUrl,
		cfg.QBittorrent.Username,
		cfg.QBittorrent.Password,
		cfg.QBittorrent.CookieDir,
	)
	if err != nil {
		logger.Fatal(err)
	}

	torrentInfo, err := torrentApi.GetTorrentInfo(torrentHash)
	if err != nil {
		apierr.PrintError(logger, err)
	}

	torrentFiles, err := torrentApi.GetTorrentFiles(torrentHash)
	if err != nil {
		apierr.PrintError(logger, err)
	}

	movieApi, err := themoviedb.NewApi(
		logger,
		cfg.MovieDb.ApiKey,
	)
	if err != nil {
		logger.Fatal(err)
	}

	episodes, err := movieApi.GetSeasonEpisodes(tvShowID, season, themoviedb.LanguageEN)
	if err != nil {
		apierr.PrintError(logger, err)
	}
	for _, episode := range episodes {
		fmt.Printf("(%d) %s (%s)\n", episode.EpisodeNumber, episode.Name, episode.AirDate)
	}

	prepareTVShow := prepare.NewService()

	params, err := mapToPrepareTvShowPrams(torrentInfo.SavePath, torrentInfo.ContentPath, episodes, torrentFiles)
	if err != nil {
		logger.Fatal(err)
	}

	result, err := prepareTVShow.PrepareTvShowSeason(params)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(result)
}
