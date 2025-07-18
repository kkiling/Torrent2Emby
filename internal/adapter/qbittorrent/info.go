package qbittorrent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

type rawTorrentInfo struct {
	AddedOn      int64   `json:"added_on"`
	AmountLeft   int64   `json:"amount_left"`
	Category     string  `json:"category"`
	Completed    int64   `json:"completed"`
	CompletionOn int64   `json:"completion_on"`
	SavePath     string  `json:"save_path"`
	ContentPath  string  `json:"content_path"`
	DlSpeed      int64   `json:"dlspeed"`
	Downloaded   int64   `json:"downloaded"`
	Eta          int     `json:"eta"`
	Hash         string  `json:"hash"`
	Name         string  `json:"name"`
	Progress     float64 `json:"progress"`
	Size         int64   `json:"size"`
	State        string  `json:"statemachine"`
	Tags         string  `json:"tags"`
	TotalSize    int64   `json:"total_size"`
	UpSpeed      int64   `json:"upspeed"`
	Uploaded     int64   `json:"uploaded"`
}

func (api *Api) GetTorrentInfo(hash string) (*TorrentInfo, error) {
	if err := api.login(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	query := url.Values{}
	query.Add("hashes", hash)

	reqURL := api.baseAPIUrl.String() + "/api/v2/torrents/info?" + query.Encode()
	resp, err := api.httpClient.Get(reqURL)
	if err != nil {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	var rawInfoList []rawTorrentInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawInfoList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(rawInfoList) == 0 {
		return nil, apierr.ContentNotFound
	}

	raw := rawInfoList[0]
	info := &TorrentInfo{
		Hash:         raw.Hash,
		Name:         raw.Name,
		Category:     raw.Category,
		Tags:         raw.Tags,
		SavePath:     raw.SavePath,
		ContentPath:  raw.ContentPath,
		State:        TorrentState(raw.State),
		AddedOn:      time.Unix(raw.AddedOn, 0),
		CompletionOn: time.Unix(raw.CompletionOn, 0),
		Eta:          raw.Eta,
		AmountLeft:   raw.AmountLeft,
		Completed:    raw.Completed,
		Downloaded:   raw.Downloaded,
		Uploaded:     raw.Uploaded,
		Size:         raw.Size,
		TotalSize:    raw.TotalSize,
		Progress:     raw.Progress,
		DlSpeed:      raw.DlSpeed,
		UpSpeed:      raw.UpSpeed,
	}

	return info, nil
}
