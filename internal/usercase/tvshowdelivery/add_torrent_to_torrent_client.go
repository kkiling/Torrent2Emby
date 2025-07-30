package tvshowdelivery

import (
	"context"
	"fmt"

	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
)

type AddTorrentParams struct {
	TVShowID TVShowID
	Magnet   string
}

// AddTorrentToTorrentClient добавление торрент раздачи в торрент клиент
func (s *Service) AddTorrentToTorrentClient(_ context.Context, params AddTorrentParams) error {
	// Создание раздачи в торрент клиенте, выставление его сразу в паузу
	err := s.torrentClient.AddTorrent(qbittorrent.TorrentAddOptions{
		Magnet:   params.Magnet,
		SavePath: s.config.TVShowTorrentSavePath,
		Category: "tvshow",
		Tags: []string{
			fmt.Sprintf("tvshowID:%d", params.TVShowID.ID),
			fmt.Sprintf("seasonNumber:%d", params.TVShowID.SeasonNumber),
		},
		Paused: false,
	})
	if err != nil {
		return fmt.Errorf("torrentClient.AddTorrent: %w", err)
	}

	return nil
}
