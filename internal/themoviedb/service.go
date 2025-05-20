package themoviedb

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kkiling/torrent2emby/internal/log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	baseApiUrl = "https://api.themoviedb.org/3"
	baseImgUrl = "https://image.tmdb.org/t/p"
	maxPages   = 5
)

type API struct {
	apiKey     string
	baseAPIUrl *url.URL
	baseImgUrl *url.URL
	httpClient *http.Client
	logger     log.Logger
	validate   *validator.Validate
}

func NewAPI(logger log.Logger, apiKey string) (*API, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	urlApi, err := url.Parse(baseApiUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}
	urlImg, err := url.Parse(baseImgUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	return &API{
		apiKey:     apiKey,
		baseAPIUrl: urlApi,
		baseImgUrl: urlImg,
		httpClient: &http.Client{Jar: jar},
		logger:     logger.Named("the_movie_db"),
		validate:   validator.New(),
	}, nil
}
