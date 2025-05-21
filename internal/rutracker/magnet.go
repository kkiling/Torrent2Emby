package rutracker

import (
	"errors"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/apierr"
	"net/http"
	"regexp"
	"strings"
)

func (api *Api) GetMagnetLink(torrentUrl string) (*MagnetInfo, error) {
	if err := api.login(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	api.logger.Debugf("Get magnet link: %s", torrentUrl)

	if !strings.HasPrefix(torrentUrl, api.baseAPIUrl.String()) {
		return nil, errors.New("invalid topic URL")
	}

	resp, err := api.httpClient.Get(torrentUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get magnet link: %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	doc, err := readerDocument(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	magnetLink := doc.Find("a.magnet-link").First()
	if magnetLink.Length() == 0 {
		return nil, errors.New("magnet link not found on the page")
	}

	magnet := magnetLink.AttrOr("href", "")
	if magnet == "" {
		return nil, errors.New("magnet link is empty")
	}

	hash := magnetLink.AttrOr("title", "")
	if hash == "" {
		re := regexp.MustCompile(`btih:([A-F0-9]{40})`)
		if match := re.FindStringSubmatch(magnet); len(match) > 1 {
			hash = match[1]
		}
	}

	return &MagnetInfo{
		Magnet: magnet,
		Hash:   strings.ToUpper(hash),
	}, nil
}
