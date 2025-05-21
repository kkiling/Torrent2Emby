package qbittorrent

import (
	"encoding/json"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/apierr"
	"net/http"
	"net/url"
)

type rawTorrentFile struct {
	Name     string  `json:"name"`
	Progress float64 `json:"progress"`
	Size     int64   `json:"size"`
}

// GetTorrentFiles возвращает список файлов в торренте по его хешу
func (api *Api) GetTorrentFiles(hash string) ([]TorrentFile, error) {
	if err := api.login(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	query := url.Values{}
	query.Add("hash", hash)

	reqURL := api.baseAPIUrl.String() + "/api/v2/torrents/files?" + query.Encode()
	resp, err := api.httpClient.Get(reqURL)
	if err != nil {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	var rawFiles []rawTorrentFile
	if err := json.NewDecoder(resp.Body).Decode(&rawFiles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	files := make([]TorrentFile, len(rawFiles))
	for i, raw := range rawFiles {
		files[i] = TorrentFile{
			Name:     raw.Name,
			Progress: raw.Progress,
			Size:     raw.Size,
		}
	}

	return files, nil
}
