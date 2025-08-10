package qbittorrent

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

func (api *Api) AddTorrent(opts TorrentAddOptions) error {
	if err := api.login(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	form := url.Values{}
	form.Set("urls", opts.Magnet)
	form.Set("savepath", opts.SavePath)

	if opts.Category != "" {
		form.Set("category", opts.Category)
	}
	if len(opts.Tags) > 0 {
		form.Set("tags", strings.Join(opts.Tags, ","))
	}
	if opts.Paused {
		form.Set("paused", "true")
	}

	postUrl := api.baseAPIUrl.String() + "/api/v2/torrents/add"
	resp, err := api.httpClient.PostForm(postUrl, form)
	if err != nil {
		return apierr.HandleStatusCodeError(api.logger, resp)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			if errRemove := api.removeCookies(); errRemove != nil {
				api.logger.Errorf("failed to remove cookies: %v", errRemove)
			}
		}
		return apierr.HandleStatusCodeError(api.logger, resp)
	}

	return nil
}
