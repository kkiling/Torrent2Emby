package emby

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/kkiling/goplatform/log"
)

type API struct {
	apiKey     string
	baseAPIUrl *url.URL
	httpClient *http.Client
	logger     log.Logger
	validate   *validator.Validate
}

func NewApi(apiKey string, baseURL string, logger log.Logger) (*API, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	urlApi, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	return &API{
		apiKey:     apiKey,
		baseAPIUrl: urlApi,
		httpClient: &http.Client{Jar: jar},
		logger:     logger.Named("emby_api"),
		validate:   validator.New(),
	}, nil
}
