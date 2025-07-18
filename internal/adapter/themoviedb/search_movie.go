package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

func (api *API) searchMovie(params SearchQuery) (*MovieSearchResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("query", params.Query)
	queryParams.Add("language", string(params.Language))
	queryParams.Add("page", fmt.Sprintf("%d", params.Page))

	getUrl := fmt.Sprintf("%s/tvshowlibrary/movie?%s", api.baseAPIUrl.String(), queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to tvshowlibrary movies: %w", apierr.HandleRequestError(api.logger, err))
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
		Page         int `json:"page"`
		TotalResults int `json:"total_results"`
		Results      []struct {
			ID            uint64  `json:"id"`
			Title         string  `json:"title"`
			OriginalTitle string  `json:"original_title"`
			Overview      string  `json:"overview"`
			PosterPath    string  `json:"poster_path"`
			ReleaseDate   string  `json:"release_date"`
			VoteAverage   float64 `json:"vote_average"`
			VoteCount     int     `json:"vote_count"`
			Popularity    float64 `json:"popularity"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	movies := make([]MovieShort, len(result.Results))
	for i, item := range result.Results {
		movies[i] = MovieShort{
			ID:            item.ID,
			Title:         item.Title,
			OriginalTitle: item.OriginalTitle,
			Overview:      item.Overview,
			PosterPath:    api.getImage(item.PosterPath),
			ReleaseDate:   parseDate(item.ReleaseDate),
			VoteAverage:   item.VoteAverage,
			VoteCount:     item.VoteCount,
			Popularity:    item.Popularity,
		}
	}

	return &MovieSearchResponse{
		Page:         result.Page,
		TotalResults: result.TotalResults,
		Results:      movies,
	}, nil
}

// SearchMovie searches for movies with sorting by popularity
func (api *API) SearchMovie(params SearchQuery) (*MovieSearchResponse, error) {
	if err := api.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("invalid tvshowlibrary query: %w", err)
	}

	page := 1
	// Get first page to determine total pages
	initialResponse, err := api.searchMovie(SearchQuery{
		Query:    params.Query,
		Language: params.Language,
		Page:     page,
		PerPage:  params.PerPage,
	})

	if err != nil {
		return nil, fmt.Errorf("searchMovie: %w", err)
	}

	allMovies := initialResponse.Results

	for page < maxPages {
		page++
		pageResponse, err := api.searchMovie(SearchQuery{
			Query:    params.Query,
			Language: params.Language,
			Page:     page,
			PerPage:  params.PerPage,
		})

		if err != nil {
			return nil, fmt.Errorf("searchMovie: %w", err)
		}

		if len(pageResponse.Results) == 0 {
			break
		}

		allMovies = append(allMovies, pageResponse.Results...)
	}

	sort.Slice(allMovies, func(i, j int) bool {
		return allMovies[i].Popularity > allMovies[j].Popularity
	})

	// Calculate pagination indices
	startIdx := (params.Page - 1) * params.PerPage
	endIdx := startIdx + params.PerPage
	if endIdx > len(allMovies) {
		endIdx = len(allMovies)
	}

	return &MovieSearchResponse{
		Page:         params.Page,
		TotalResults: len(allMovies),
		Results:      allMovies[startIdx:endIdx],
	}, nil
}
