package rutracker

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/kkiling/goplatform/log"
)

const (
	cookeFile = "rutracker_cookies.gob"
	apiUrl    = "https://rutracker.org/forum/"
)

type Api struct {
	username   string
	password   string
	cookiesDir string
	baseAPIUrl *url.URL
	httpClient *http.Client
	logger     log.Logger
}

func NewApi(logger log.Logger, username, password, cookiesDir string) (*Api, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	baseAPIUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	return &Api{
		username:   username,
		password:   password,
		cookiesDir: cookiesDir,
		baseAPIUrl: baseAPIUrl,
		httpClient: &http.Client{Jar: jar},
		logger:     logger.Named("rutracker"),
	}, nil
}
