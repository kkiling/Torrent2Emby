package contentdelivery

import (
	"context"
	"fmt"

	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
)

type WaitingTorrentDownloadCompleteParams struct {
	Hash string
}

// WaitingTorrentDownloadComplete ожидание завершения окончания скачивания раздачи
func (s *Service) WaitingTorrentDownloadComplete(_ context.Context, params WaitingTorrentDownloadCompleteParams) (*TorrentDownloadStatus, error) {
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
	case qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP,
		qbittorrent.TorrentStateStalledUP:
		return &TorrentDownloadStatus{
			Progress:   torrentInfo.Progress,
			IsComplete: true,
		}, nil
	}

	return &TorrentDownloadStatus{
		Progress:   torrentInfo.Progress,
		IsComplete: false,
	}, nil
}
