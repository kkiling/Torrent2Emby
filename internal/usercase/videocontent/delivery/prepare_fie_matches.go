package delivery

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/adapter/matchtvshow"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
)

type PreparingFileMatchesParams struct {
	Hash     string
	TVShowID common.TVShowID
}

func mapFile(file matchtvshow.TorrentFile) FileInfo {
	return FileInfo{
		RelativePath: file.RelativePath,
		FullPath:     file.FullPath,
		Size:         file.Size,
		Extension:    file.Extension,
	}
}

func mapTrack(tracks []matchtvshow.PrepareTrack) []Track {
	return lo.Map(tracks, func(item matchtvshow.PrepareTrack, index int) Track {
		return Track{
			Name:     item.Name,
			Language: item.Language,
			File:     mapFile(item.File),
		}
	})
}

func mapContentMatchesFromPrepareTVShowSeason(
	prepareResult *matchtvshow.PrepareTVShowSeason,
) ([]ContentMatches, error) {
	result := make([]ContentMatches, 0, len(prepareResult.Episodes))
	for _, prepareEpisode := range prepareResult.Episodes {
		content := ContentMatches{
			Episode: EpisodeInfo{
				SeasonNumber:  prepareEpisode.Episode.SeasonNumber,
				EpisodeName:   prepareEpisode.Episode.EpisodeName,
				EpisodeNumber: prepareEpisode.Episode.EpisodeNumber,
			},
			Video: VideoFile{
				File: mapFile(prepareEpisode.VideoFile.File),
			},
			AudioFiles: mapTrack(prepareEpisode.AudioFiles),
			Subtitles:  mapTrack(prepareEpisode.Subtitles),
		}

		result = append(result, content)
	}

	return result, nil
}

func mapToPrepareTvShowPrams(
	basePath, savePath, contentPath string,
	episodes []tvshowlibrary.Episode,
	torrentFiles []qbittorrent.TorrentFile,
) (*matchtvshow.PrepareTvShowPrams, error) {
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
	var prepareTorrentFiles []matchtvshow.TorrentFile
	for _, file := range torrentFiles {
		relFile, err := filepath.Rel(relPath, file.Name)
		if err != nil {
			return nil, fmt.Errorf("filepath.Rel: %w", err)
		}
		prepareTorrentFiles = append(prepareTorrentFiles, matchtvshow.TorrentFile{
			RelativePath: relFile,
			FullPath:     filepath.Join(fullPath, relFile),
			Extension:    strings.ToLower(filepath.Ext(relFile)),
			Size:         file.Size,
		})
	}

	return &matchtvshow.PrepareTvShowPrams{
		Episodes: lo.Map(episodes, func(episode tvshowlibrary.Episode, _ int) matchtvshow.EpisodeInfo {
			return matchtvshow.EpisodeInfo{
				SeasonNumber:  episode.EpisodeNumber,
				EpisodeNumber: episode.EpisodeNumber,
				EpisodeName:   episode.Name,
			}
		}),
		TorrentFiles: prepareTorrentFiles,
	}, nil
}

// PrepareFileMatches получение информации о файлах раздачи
func (s *Service) PrepareFileMatches(ctx context.Context, params PreparingFileMatchesParams) ([]ContentMatches, error) {
	// Достаем инфу о торрент раздаче
	torrentInfo, err := s.torrentClient.GetTorrentInfo(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if torrentInfo == nil {
		return nil, fmt.Errorf("torrentInfo not found: %w", ucerr.NotFound)
	}

	switch torrentInfo.State {
	case qbittorrent.TorrentStatePausedDL, qbittorrent.TorrentStateStoppedDL:
		if err := s.torrentClient.ResumeTorrent(params.Hash); err != nil {
			return nil, fmt.Errorf("torrentClient.ResumeTorrent: %w", err)
		}
	case
		qbittorrent.TorrentStateDownloading,
		qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP,
		qbittorrent.TorrentStateStalledUP:
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
		TVShowID:     params.TVShowID.ID,
		SeasonNumber: params.TVShowID.SeasonNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetSeasonEpisodes: %w", err)
	}

	// Подготавливаем параметры для преобразования файлов
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

	// Получаем инфу о метчах
	result, err := mapContentMatchesFromPrepareTVShowSeason(prepareResult)
	if err != nil {
		return nil, fmt.Errorf("mapContentMatchesFromPrepareTVShowSeason: %w", err)
	}

	return result, nil
}
