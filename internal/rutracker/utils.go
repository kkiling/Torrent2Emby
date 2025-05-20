package rutracker

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kkiling/torrent2emby/internal/log"
	"golang.org/x/net/html/charset"
	"io"
	"net"
	"net/http"
	"strings"
)

func readerDocument(body io.Reader) (*goquery.Document, error) {
	utf8Reader, err := charset.NewReader(body, "text/html")
	if err != nil {
		return nil, fmt.Errorf("failed to create charset reader: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}

	return doc, nil
}

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
	if strings.Contains(err.Error(), "TLS handshake timeout") {
		return ServiceUnavailableErr
	}

	return err
}

func emptyTorrentResponse() *TorrentResponse {
	return &TorrentResponse{
		Results:      []Torrent{},
		Page:         1,
		TotalResults: 1,
	}
}
