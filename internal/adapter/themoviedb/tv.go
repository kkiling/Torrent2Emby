package themoviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

func (api *API) GetTV(ctx context.Context, tvID uint64, language Language) (*TVShow, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("language", string(language))

	getUrl := fmt.Sprintf("%s/tv/%d?%s", api.baseAPIUrl.String(), tvID, queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get tv info: %w", apierr.HandleRequestError(api.logger, err))
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
		FirstAirDate string `json:"first_air_date"`
		Genres       []struct {
			Name string `json:"name"`
		} `json:"genres"`
		ID               uint64   `json:"id"`
		LastAirDate      string   `json:"last_air_date"`
		Name             string   `json:"name"`
		NextEpisodeToAir string   `json:"next_episode_to_air"`
		NumberOfEpisodes uint32   `json:"number_of_episodes"`
		NumberOfSeasons  uint32   `json:"number_of_seasons"`
		OriginCountry    []string `json:"origin_country"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		Poster           string   `json:"poster_path"`
		Seasons          []struct {
			AirDate      string  `json:"air_date"`
			EpisodeCount uint32  `json:"episode_count"`
			ID           uint64  `json:"id"`
			Name         string  `json:"name"`
			Overview     string  `json:"overview"`
			Poster       string  `json:"poster_path"`
			SeasonNumber uint8   `json:"season_number"`
			VoteAverage  float64 `json:"vote_average"`
		} `json:"seasons"`
		Status      string  `json:"status"`
		Tagline     string  `json:"tagline"`
		Type        string  `json:"type"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   uint32  `json:"vote_count"`
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
			AirDate:      parseDate(s.AirDate),
			EpisodeCount: s.EpisodeCount,
			ID:           s.ID,
			Name:         s.Name,
			Overview:     s.Overview,
			Poster:       api.getImage(s.Poster),
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
			Poster:       api.getImage(result.Poster),
			FirstAirDate: parseDate(result.FirstAirDate),
			VoteAverage:  result.VoteAverage,
			VoteCount:    result.VoteCount,
			Popularity:   result.Popularity,
		},
		Backdrop:         api.getImage(result.BackdropPath),
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
