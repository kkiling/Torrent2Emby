package qbittorrent

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/kkiling/goplatform/log"
)

const (
	cookeFile = "qbittorrentr_cookies.gob"
)

// Api представляет клиент для работы с API qBittorrent
type Api struct {
	baseAPIUrl *url.URL
	username   string
	password   string
	cookiesDir string
	httpClient *http.Client
	logger     log.Logger
}

func NewApi(logger log.Logger, baseURL, username, password, cookiesDir string) (*Api, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	baseAPIUrl, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	return &Api{
		baseAPIUrl: baseAPIUrl,
		username:   username,
		password:   password,
		cookiesDir: cookiesDir,
		httpClient: &http.Client{Jar: jar},
		logger:     logger.Named("qbittorrent"),
	}, nil
}
