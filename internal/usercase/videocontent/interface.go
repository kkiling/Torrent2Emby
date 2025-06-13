package videocontent

import (
	"context"
	"github.com/kkiling/torrent2emby/internal/adapter/qbittorrent"
	"github.com/kkiling/torrent2emby/internal/adapter/rutracker"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

type Repository interface {
	// GetVideoContentByMediaID Получение видеоконтента для сезона сериала или фильма
	GetVideoContentByMediaID(ctx context.Context, id MediaID) ([]VideoContent, error)
	// SaveVideoContent сохранение информации о видео контенте
	SaveVideoContent(ctx context.Context, content *VideoContent) error
	// Сохранение разультатов поиска раздач
}

type TVShowLibrary interface {
	GetTVShowInfo(ctx context.Context, params tvshowlibrary.GetTVShowParams) (*tvshowlibrary.GetTVShowResult, error)
}

type TorrentSite interface {
	SearchTorrents(query string) (*rutracker.TorrentResponse, error)
	GetMagnetLink(torrentUrl string) (*rutracker.MagnetInfo, error)
}

type TorrentClient interface {
	AddTorrent(opts qbittorrent.TorrentAddOptions) error
}
