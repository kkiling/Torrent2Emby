package qbittorrent

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

func (api *Api) ResumeTorrent(hash string) error {
	if err := api.login(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	form := url.Values{}
	form.Set("hashes", hash)

	postUrl := api.baseAPIUrl.String() + "/api/v2/torrents/start"
	resp, err := api.httpClient.PostForm(postUrl, form)
	if err != nil {
		return apierr.HandleStatusCodeError(api.logger, resp)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apierr.HandleStatusCodeError(api.logger, resp)
	}

	return nil
}
