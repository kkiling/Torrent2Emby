package tvshowlibrary

import (
	themoviedb2 "github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/samber/lo"
)

func mapImage(image *themoviedb2.Image) *Image {
	return &Image{
		W92:      image.W92,
		W154:     image.W154,
		W185:     image.W185,
		W342:     image.W342,
		W500:     image.W500,
		W780:     image.W780,
		Original: image.Original,
	}
}

func mapTVShowShort(item themoviedb2.TVShowShort) TVShowShort {
	return TVShowShort{
		ID:           item.ID,
		Name:         item.Name,
		OriginalName: item.OriginalName,
		Overview:     item.Overview,
		PosterPath:   mapImage(item.PosterPath),
		FirstAirDate: item.FirstAirDate,
		VoteAverage:  item.VoteAverage,
		VoteCount:    item.VoteCount,
		Popularity:   item.Popularity,
	}
}

func mapSeasonShort(season themoviedb2.Season) SeasonShort {
	return SeasonShort{
		ID:           season.ID,
		AirDate:      season.AirDate,
		EpisodeCount: season.EpisodeCount,
		Name:         season.Name,
		Overview:     season.Overview,
		PosterPath:   mapImage(season.PosterPath),
		SeasonNumber: season.SeasonNumber,
		VoteAverage:  season.VoteAverage,
	}
}

func mapTvShowSearchResult(response *themoviedb2.TVShowSearchResponse) (*TVShowSearchResult, error) {
	return &TVShowSearchResult{Results: lo.Map(response.Results, func(item themoviedb2.TVShowShort, index int) TVShowShort {
		return mapTVShowShort(item)
	})}, nil
}

func mapGetTVShowResult(response *themoviedb2.TVShow) (*GetTVShowResult, error) {
	return &GetTVShowResult{
		TVShowInfo: TVShowSearch{
			TVShow: TVShow{
				TVShowShort:      mapTVShowShort(response.TVShowShort),
				BackdropPath:     mapImage(response.BackdropPath),
				Genres:           response.Genres,
				LastAirDate:      response.LastAirDate,
				NextEpisodeToAir: response.NextEpisodeToAir,
				NumberOfEpisodes: response.NumberOfEpisodes,
				NumberOfSeasons:  response.NumberOfSeasons,
				OriginCountry:    response.OriginCountry,

				Status:  response.Status,
				Tagline: response.Tagline,
				Type:    response.Type,
			},
			Seasons: lo.Map(response.Seasons, func(item themoviedb2.Season, index int) SeasonShort {
				return mapSeasonShort(item)
			}),
		},
	}, nil
}

func mapGetSeasonEpisodesResult(response []themoviedb2.Episode) (*GetSeasonEpisodesResult, error) {
	return &GetSeasonEpisodesResult{
		Episodes: lo.Map(response, func(item themoviedb2.Episode, index int) Episode {
			return Episode{
				ID:            item.ID,
				AirDate:       item.AirDate,
				EpisodeNumber: item.EpisodeNumber,
				EpisodeType:   item.EpisodeType,
				Name:          item.Name,
				Overview:      item.Overview,
				Runtime:       item.Runtime,
				StillPath:     mapImage(item.StillPath),
				VoteAverage:   item.VoteAverage,
				VoteCount:     item.VoteCount,
			}
		}),
	}, nil
}
