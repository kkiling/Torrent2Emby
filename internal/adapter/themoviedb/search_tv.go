package themoviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

func (api *API) searchTV(params SearchQuery) (*TVShowSearchResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("query", params.Query)
	queryParams.Add("language", string(params.Language))
	queryParams.Add("page", fmt.Sprintf("%d", params.Page))

	getUrl := fmt.Sprintf("%s/search/tv?%s", api.baseAPIUrl.String(), queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to tvshowlibrary tv: %w", apierr.HandleRequestError(api.logger, err))
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
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
		Results      []struct {
			ID           uint64  `json:"id"`
			Name         string  `json:"name"`
			OriginalName string  `json:"original_name"`
			Overview     string  `json:"overview"`
			Poster       string  `json:"poster_path"`
			FirstAirDate string  `json:"first_air_date"`
			VoteAverage  float64 `json:"vote_average"`
			VoteCount    uint32  `json:"vote_count"`
			Popularity   float64 `json:"popularity"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	movies := make([]TVShowShort, len(result.Results))
	for i, item := range result.Results {
		movies[i] = TVShowShort{
			ID:           item.ID,
			Name:         item.Name,
			OriginalName: item.OriginalName,
			Overview:     item.Overview,
			Poster:       api.getImage(item.Poster),
			FirstAirDate: parseDate(item.FirstAirDate),
			VoteAverage:  item.VoteAverage,
			VoteCount:    item.VoteCount,
			Popularity:   item.Popularity,
		}
	}

	return &TVShowSearchResponse{
		Page:         result.Page,
		TotalResults: result.TotalResults,
		Results:      movies,
	}, nil
}

func (api *API) SearchTV(ctx context.Context, params SearchQuery) (*TVShowSearchResponse, error) {
	if err := api.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("invalid tvshowlibrary query: %w", err)
	}

	page := 1
	// Get first page to determine total pages
	initialResponse, err := api.searchTV(SearchQuery{
		Query:    params.Query,
		Language: params.Language,
		Page:     page,
		PerPage:  params.PerPage,
	})

	if err != nil {
		return nil, fmt.Errorf("searchTV: %w", err)
	}

	allTVs := initialResponse.Results

	for page < maxPages {
		page++
		pageResponse, err := api.searchTV(SearchQuery{
			Query:    params.Query,
			Language: params.Language,
			Page:     page,
			PerPage:  params.PerPage,
		})

		if err != nil {
			return nil, fmt.Errorf("searchTV: %w", err)
		}

		if len(pageResponse.Results) == 0 {
			break
		}

		allTVs = append(allTVs, pageResponse.Results...)
	}

	sort.Slice(allTVs, func(i, j int) bool {
		return allTVs[i].Popularity > allTVs[j].Popularity
	})

	// Calculate pagination indices
	startIdx := (params.Page - 1) * params.PerPage
	endIdx := startIdx + params.PerPage
	if endIdx > len(allTVs) {
		endIdx = len(allTVs)
	}

	return &TVShowSearchResponse{
		Page:         params.Page,
		TotalResults: len(allTVs),
		Results:      allTVs[startIdx:endIdx],
	}, nil
}
