package qbittorrent

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
	"net/http"
	"net/url"
)

func (api *Api) DeleteTorrent(hash string, deleteFiles bool) error {
	if err := api.login(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	form := url.Values{}
	form.Set("hashes", hash)
	if deleteFiles {
		form.Set("deleteFiles", "true")
	}

	postUrl := api.baseAPIUrl.String() + "/api/v2/torrents/delete"
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
