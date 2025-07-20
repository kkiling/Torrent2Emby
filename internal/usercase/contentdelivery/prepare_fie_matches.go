package contentdelivery

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/adapter/prepare/tvshow"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

type PreparingFileMatchesParams struct {
	Hash    string
	MediaID MediaID
}

func mapFile(file tvshow.TorrentFile) FileInfo {
	return FileInfo{
		RelativePath: file.RelativePath,
		FullPath:     file.FullPath,
		Size:         file.Size,
		Extension:    file.Extension,
	}
}

func mapTrack(tracks []tvshow.PrepareTrack) []Track {
	return lo.Map(tracks, func(item tvshow.PrepareTrack, index int) Track {
		return Track{
			Name:     item.Name,
			Language: item.Language,
			File:     mapFile(item.File),
		}
	})
}

func mapContentMatchesFromPrepareTVShowSeason(
	prepareResult *tvshow.PrepareTVShowSeason,
	episodes []tvshowlibrary.Episode,
) []ContentMatches {
	result := make([]ContentMatches, len(prepareResult.Episodes))
	for _, episode := range prepareResult.Episodes {
		ep := episodes[episode.EpisodeNumber]

		items := ContentMatches{
			ContentInfo: ContentInfo{
				Name: fmt.Sprintf("%d %s", ep.EpisodeNumber, ep.Name),
			},
			Video: VideoFile{
				File: mapFile(episode.VideoFile.File),
			},
			AudioFiles: mapTrack(episode.AudioFiles),
			Subtitles:  mapTrack(episode.Subtitles),
		}

		result = append(result, items)
	}

	return result
}

func mapToPrepareTvShowPrams(
	basePath, savePath, contentPath string,
	episodes []tvshowlibrary.Episode,
	torrentFiles []qbittorrent.TorrentFile,
) (*tvshow.PrepareTvShowPrams, error) {
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
	var prepareTorrentFiles []tvshow.TorrentFile
	for _, file := range torrentFiles {
		relFile, err := filepath.Rel(relPath, file.Name)
		if err != nil {
			return nil, fmt.Errorf("filepath.Rel: %w", err)
		}
		prepareTorrentFiles = append(prepareTorrentFiles, tvshow.TorrentFile{
			RelativePath: relFile,
			FullPath:     filepath.Join(fullPath, relFile),
			Extension:    strings.ToLower(filepath.Ext(relFile)),
			Size:         file.Size,
		})
	}

	return &tvshow.PrepareTvShowPrams{
		Episodes: lo.Map(episodes, func(episode tvshowlibrary.Episode, _ int) tvshow.Episode {
			return tvshow.Episode{
				EpisodeNumber: episode.EpisodeNumber,
			}
		}),
		TorrentFiles: prepareTorrentFiles,
	}, nil
}

// PrepareFileMatches получение информации о файлах раздачи
func (s *Service) PrepareFileMatches(ctx context.Context, params PreparingFileMatchesParams) ([]ContentMatches, error) {
	if params.MediaID.MovieID != nil && params.MediaID.TVShow == nil {
		return nil, fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	// Достаем инфу о торрент раздаче
	torrentInfo, err := s.torrentClient.GetTorrentInfo(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if torrentInfo == nil {
		return nil, fmt.Errorf("torrentInfo not found: %w", ucerr.NotFound)
	}

	switch torrentInfo.State {
	case qbittorrent.TorrentStateDownloading,
		qbittorrent.TorrentStatePausedDL,
		qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP:
		// Файлы начали скачиваться, значит можем получить информацию о файлах
	default:
		// Ошибки как таковой нет, придем в следующий раз
		return nil, nil
	}

	torrentFiles, err := s.torrentClient.GetTorrentFiles(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if len(torrentFiles) == 0 {
		return nil, fmt.Errorf("torrentFiles not found: %w", ucerr.NotFound)
	}

	// Достаем инфу о эпизодах
	episodes, err := s.tvShowLibrary.GetSeasonEpisodes(ctx, tvshowlibrary.GetSeasonEpisodesParams{
		TVShowID:     params.MediaID.TVShow.TVShowID,
		SeasonNumber: params.MediaID.TVShow.SeasonNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetSeasonEpisodes: %w", err)
	}

	// Подготавливаем параметры для преобразования файлов
	// TODO: подумать так как торрент работает в контейнере у него абсолютный файл начинается с /download
	prepareParams, err := mapToPrepareTvShowPrams(
		s.config.BasePath,
		torrentInfo.SavePath,
		torrentInfo.ContentPath,
		episodes.Items,
		torrentFiles,
	)

	if err != nil {
		return nil, fmt.Errorf("mapToPrepareTvShowPrams: %w", err)
	}

	prepareResult, err := s.prepareTVShow.PrepareTvShowSeason(prepareParams)
	if err != nil {
		return nil, fmt.Errorf("prepareTVShow.PrepareTvShowSeason: %w", err)
	}

	result := mapContentMatchesFromPrepareTVShowSeason(prepareResult, episodes.Items)

	return result, nil
}
