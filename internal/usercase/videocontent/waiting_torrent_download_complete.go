package videocontent

import (
	"context"
	"fmt"

	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
)

type WaitingTorrentDownloadCompleteParams struct {
	Hash string
}

// waitingTorrentDownloadComplete ожидание завершения окончания скачивания раздачи
func (s *Service) waitingTorrentDownloadComplete(_ context.Context, params WaitingTorrentDownloadCompleteParams) (*TorrentDownloadStatus, error) {
	// Достаем инфу о торрент раздаче
	torrentInfo, err := s.torrentClient.GetTorrentInfo(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if torrentInfo == nil {
		return nil, fmt.Errorf("torrentInfo not found: %w", ucerr.NotFound)
	}

	switch torrentInfo.State {
	case qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP:
		return &TorrentDownloadStatus{
			Progress:   torrentInfo.Progress,
			IsComplete: true,
		}, nil
	default:
		return &TorrentDownloadStatus{
			Progress:   torrentInfo.Progress,
			IsComplete: false,
		}, nil
	}
}
