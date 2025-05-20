package themoviedb

import (
	"errors"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/log"
	"net"
	"net/http"
	"strings"
	"time"
)

func handleStatusCodeError(log log.Logger, resp *http.Response) error {
	log.Errorf("Response status %s: %s", resp.Request.URL.Path, resp.Status)
	if resp.StatusCode == http.StatusUnauthorized {
		return NotAuthorizedErr
	} else if resp.StatusCode == 522 || resp.StatusCode == 521 {
		return ServiceUnavailableErr
	}
	return fmt.Errorf("search failed with status code: %d", resp.StatusCode)
}

func handleRequestError(log log.Logger, err error) error {
	log.Errorf("Request failed: %v", err)
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ServiceUnavailableErr
	}
	if strings.Contains(err.Error(), "connection reset by peer") {
		return ServiceUnavailableErr
	}
	if strings.Contains(err.Error(), "connection refused") {
		return ServiceUnavailableErr
	}
	if strings.Contains(err.Error(), "TLS handshake timeout") {
		return ServiceUnavailableErr
	}

	return err
}

func (api *API) getImage(path string) *Image {
	if path == "" {
		return nil
	}
	urlImg := api.baseImgUrl.String()
	return &Image{
		W92:      fmt.Sprintf("%s/w92%s", urlImg, path),
		W154:     fmt.Sprintf("%s/w154%s", urlImg, path),
		W185:     fmt.Sprintf("%s/w185%s", urlImg, path),
		W342:     fmt.Sprintf("%s/w342%s", urlImg, path),
		W500:     fmt.Sprintf("%s/w500%s", urlImg, path),
		W780:     fmt.Sprintf("%s/w780%s", urlImg, path),
		Original: fmt.Sprintf("%s/original%s", urlImg, path),
	}
}

func parseDate(dateStr string) time.Time {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
