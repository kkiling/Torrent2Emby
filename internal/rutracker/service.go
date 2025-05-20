package rutracker

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	cookeFile = "rutracker_cookies.gob"
	apiUrl    = "https://rutracker.org/forum/"
)

type API struct {
	username   string
	password   string
	cookiesDir string
	baseAPIUrl *url.URL
	httpClient *http.Client
	logger     log.Logger
}

func NewAPI(logger log.Logger, username, password, cookiesDir string) (*API, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	baseAPIUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	return &API{
		username:   username,
		password:   password,
		cookiesDir: cookiesDir,
		baseAPIUrl: baseAPIUrl,
		httpClient: &http.Client{Jar: jar},
		logger:     logger.Named("rutracker"),
	}, nil
}
