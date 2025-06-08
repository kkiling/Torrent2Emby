package main

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/apierr"
	"github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/mkvmerge"
	prepare "github.com/kkiling/torrent2emby/internal/prepare/tvshow"
	"github.com/kkiling/torrent2emby/internal/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/themoviedb"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
)

const (
	torrentHash = "556539d480ab08a99deb5f2889146f665627f24f"
	tvShowID    = 88803 // Сага о винланде
	season      = 2
)

func mapToPrepareTvShowPrams(
	basePath, savePath, contentPath string,
	episodes []themoviedb.Episode,
	torrentFiles []qbittorrent.TorrentFile,
) (*prepare.PrepareTvShowPrams, error) {
	fullPath := filepath.Join(basePath, contentPath)
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
			FullPath:     filepath.Join(fullPath, relFile),
			Extension:    strings.ToLower(filepath.Ext(relFile)),
			Size:         file.Size,
		})
	}

	return &prepare.PrepareTvShowPrams{
		Episodes: lo.Map(episodes, func(episode themoviedb.Episode, _ int) prepare.Episode {
			return prepare.Episode{
				EpisodeNumber: episode.EpisodeNumber,
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
	movieApi, err := themoviedb.NewApi(
		logger,
		cfg.MovieDb.ApiKey,
	)
	if err != nil {
		logger.Fatal(err)
	}
	merge := mkvmerge.NewService()
	prepareTVShow := prepare.NewService(merge)
	// ****

	// Достаем инфу о торрент раздаче
	torrentInfo, err := torrentApi.GetTorrentInfo(torrentHash)
	if err != nil {
		apierr.PrintError(logger, err)
	}

	torrentFiles, err := torrentApi.GetTorrentFiles(torrentHash)
	if err != nil {
		apierr.PrintError(logger, err)
	}

	// Достаем инфу о эпизодах
	episodes, err := movieApi.GetSeasonEpisodes(tvShowID, season, themoviedb.LanguageEN)
	if err != nil {
		apierr.PrintError(logger, err)
	}

	for _, episode := range episodes {
		fmt.Printf("(%d) %s (%s)\n", episode.EpisodeNumber, episode.Name, episode.AirDate)
	}

	// Подготавливаем параметры для преобразования файлов
	// TODO: подумать так как торрент работает в контейнере у него абсолютный файл начинается с /download
	params, err := mapToPrepareTvShowPrams("/nfs", torrentInfo.SavePath, torrentInfo.ContentPath, episodes, torrentFiles)
	if err != nil {
		logger.Fatal(err)
	}

	prepareResult, err := prepareTVShow.PrepareTvShowSeason(params)
	if err != nil {
		logger.Fatal(err)
	}

	//
	epPrepare := prepareResult.Episodes[1]
	// по  ep.Episode.EpisodeNumber достаем эпизод
	epInfo := episodes[epPrepare.Episode.EpisodeNumber]

	newEpisodeFileName := fmt.Sprintf("%03d %s.%s", epInfo.EpisodeNumber, epInfo.Name, epPrepare.VideoFile.File.Extension)
	//

	mergeParams := mkvmerge.MergeParams{
		VideoInputFile:  epPrepare.VideoFile.File.FullPath,
		VideoOutputFile: filepath.Join("/home/kiling/Downloads/1", newEpisodeFileName),
		AudioTracks: lo.Map(epPrepare.AudioFiles, func(item prepare.PrepareTrack, index int) mkvmerge.Track {
			return mkvmerge.Track{
				Path:     item.File.FullPath,
				Language: item.Language,
				Name:     item.Name,
				Default:  index == 0,
			}
		}),
		SubtitleTracks: lo.Map(epPrepare.Subtitles, func(item prepare.PrepareTrack, index int) mkvmerge.Track {
			return mkvmerge.Track{
				Path:     item.File.FullPath,
				Language: item.Language,
				Name:     item.Name,
				Default:  false,
			}
		}),
	}

	err = merge.Merge(mergeParams)
	if err != nil {
		logger.Fatal(err)
	}

	info, err := merge.GetMediaInfo(mergeParams.VideoOutputFile)
	if err != nil {
		logger.Fatal(err)
	}

	// Проверяем что

	fmt.Println(info)
}
