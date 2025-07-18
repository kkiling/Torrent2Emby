package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

func (api *API) GetMovie(movieID uint64, language Language) (*Movie, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("language", string(language))

	getUrl := fmt.Sprintf("%s/movie/%d?%s", api.baseAPIUrl.String(), movieID, queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie info: %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		BackdropPath string `json:"backdrop_path"`
		Budget       int64  `json:"budget"`
		Genres       []struct {
			Name string `json:"name"`
		} `json:"genres"`
		ID               uint64   `json:"id"`
		ImdbID           string   `json:"imdb_id"`
		OriginCountry    []string `json:"origin_country"`
		OriginalLanguage string   `json:"original_language"`
		OriginalTitle    string   `json:"original_title"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		PosterPath       string   `json:"poster_path"`
		ReleaseDate      string   `json:"release_date"`
		Revenue          int64    `json:"revenue"`
		Runtime          int      `json:"runtime"`
		Status           string   `json:"status"`
		Tagline          string   `json:"tagline"`
		Title            string   `json:"title"`
		VoteAverage      float64  `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	genres := make([]string, len(result.Genres))
	for i, g := range result.Genres {
		genres[i] = g.Name
	}

	return &Movie{
		MovieShort: MovieShort{
			ID:            result.ID,
			Title:         result.Title,
			OriginalTitle: result.OriginalTitle,
			Overview:      result.Overview,
			PosterPath:    api.getImage(result.PosterPath),
			ReleaseDate:   parseDate(result.ReleaseDate),
			VoteAverage:   result.VoteAverage,
			VoteCount:     result.VoteCount,
			Popularity:    result.Popularity,
		},
		BackdropPath:     api.getImage(result.BackdropPath),
		Budget:           result.Budget,
		Genres:           genres,
		OriginCountry:    result.OriginCountry,
		OriginalLanguage: result.OriginalLanguage,
		Revenue:          result.Revenue,
		Runtime:          result.Runtime,
		Status:           result.Status,
		Tagline:          result.Tagline,
	}, nil
}
