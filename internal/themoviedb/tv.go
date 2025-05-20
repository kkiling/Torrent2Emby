package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (api *API) GetTV(tvID uint64, language Language) (*TVShow, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("language", string(language))

	getUrl := fmt.Sprintf("%s/tv/%d?%s", api.baseAPIUrl.String(), tvID, queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get tv info: %w", handleRequestError(api.logger, err))
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
		BackdropPath string `json:"backdrop_path"`
		FirstAirDate string `json:"first_air_date"`
		Genres       []struct {
			Name string `json:"name"`
		} `json:"genres"`
		ID               uint64   `json:"id"`
		LastAirDate      string   `json:"last_air_date"`
		Name             string   `json:"name"`
		NextEpisodeToAir string   `json:"next_episode_to_air"`
		NumberOfEpisodes int      `json:"number_of_episodes"`
		NumberOfSeasons  int      `json:"number_of_seasons"`
		OriginCountry    []string `json:"origin_country"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		PosterPath       string   `json:"poster_path"`
		Seasons          []struct {
			AirDate      string  `json:"air_date"`
			EpisodeCount int     `json:"episode_count"`
			ID           int     `json:"id"`
			Name         string  `json:"name"`
			Overview     string  `json:"overview"`
			PosterPath   string  `json:"poster_path"`
			SeasonNumber int     `json:"season_number"`
			VoteAverage  float64 `json:"vote_average"`
		} `json:"seasons"`
		Status      string  `json:"status"`
		Tagline     string  `json:"tagline"`
		Type        string  `json:"type"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	genres := make([]string, len(result.Genres))
	for i, g := range result.Genres {
		genres[i] = g.Name
	}

	seasons := make([]Season, len(result.Seasons))
	for i, s := range result.Seasons {
		seasons[i] = Season{
			AirDate:      s.AirDate,
			EpisodeCount: s.EpisodeCount,
			ID:           s.ID,
			Name:         s.Name,
			Overview:     s.Overview,
			PosterPath:   api.getImage(s.PosterPath),
			SeasonNumber: s.SeasonNumber,
			VoteAverage:  s.VoteAverage,
		}
	}

	return &TVShow{
		TVShowShort: TVShowShort{
			ID:           result.ID,
			Name:         result.Name,
			OriginalName: result.OriginalName,
			Overview:     result.Overview,
			PosterPath:   api.getImage(result.PosterPath),
			FirstAirDate: parseDate(result.FirstAirDate),
			VoteAverage:  result.VoteAverage,
			VoteCount:    result.VoteCount,
			Popularity:   result.Popularity,
		},
		BackdropPath:     api.getImage(result.BackdropPath),
		Genres:           genres,
		LastAirDate:      parseDate(result.LastAirDate),
		NextEpisodeToAir: parseDate(result.NextEpisodeToAir),
		NumberOfEpisodes: result.NumberOfEpisodes,
		NumberOfSeasons:  result.NumberOfSeasons,
		OriginCountry:    result.OriginCountry,
		Seasons:          seasons,
		Status:           result.Status,
		Tagline:          result.Tagline,
		Type:             result.Type,
	}, nil
}
