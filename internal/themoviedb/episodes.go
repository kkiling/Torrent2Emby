package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (api *API) GetSeasonEpisodes(tvID uint64, seasonNumber int, language Language) ([]Episode, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("language", string(language))

	getUrl := fmt.Sprintf("%s/tv/%d/season/%d?%s", api.baseAPIUrl.String(), tvID, seasonNumber, queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get season episodes: %w", handleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleStatusCodeError(api.logger, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Episodes []struct {
			AirDate       string  `json:"air_date"`
			EpisodeNumber int     `json:"episode_number"`
			EpisodeType   string  `json:"episode_type"`
			ID            int     `json:"id"`
			Name          string  `json:"name"`
			Overview      string  `json:"overview"`
			Runtime       int     `json:"runtime"`
			StillPath     string  `json:"still_path"`
			VoteAverage   float64 `json:"vote_average"`
			VoteCount     int     `json:"vote_count"`
		} `json:"episodes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	episodes := make([]Episode, len(result.Episodes))
	for i, ep := range result.Episodes {
		episodes[i] = Episode{
			AirDate:       parseDate(ep.AirDate),
			EpisodeNumber: ep.EpisodeNumber,
			EpisodeType:   ep.EpisodeType,
			ID:            ep.ID,
			Name:          ep.Name,
			Overview:      ep.Overview,
			Runtime:       ep.Runtime,
			StillPath:     api.getImage(ep.StillPath),
			VoteAverage:   ep.VoteAverage,
			VoteCount:     ep.VoteCount,
		}
	}

	return episodes, nil
}
