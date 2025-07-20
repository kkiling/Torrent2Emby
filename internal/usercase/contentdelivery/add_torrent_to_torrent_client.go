package contentdelivery

import (
	"context"
	"fmt"

	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
)

type AddTorrentParams struct {
	MediaID MediaID
	Magnet  string
}

// AddTorrentToTorrentClient добавление торрент раздачи в торрент клиент
func (s *Service) AddTorrentToTorrentClient(_ context.Context, params AddTorrentParams) error {
	if params.MediaID.MovieID != nil && params.MediaID.TVShow == nil {
		return fmt.Errorf("movie is not supported yet: %w", ucerr.InvalidArgument)
	}

	// Создание раздачи в торрент клиенте, выставление его сразу в паузу
	err := s.torrentClient.AddTorrent(qbittorrent.TorrentAddOptions{
		Magnet:   params.Magnet,
		SavePath: s.config.TVShowTorrentSavePath,
		Category: "tvshow",
		Tags: []string{
			fmt.Sprintf("tvshowID:%d", params.MediaID.TVShow.TVShowID),
			fmt.Sprintf("seasonNumber:%d", params.MediaID.TVShow.SeasonNumber),
		},
		Paused: false,
	})
	if err != nil {
		return fmt.Errorf("torrentClient.AddTorrent: %w", err)
	}

	return nil
}
