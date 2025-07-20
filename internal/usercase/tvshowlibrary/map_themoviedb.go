package tvshowlibrary

import (
	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
)

func mapImage(image *themoviedb.Image) *Image {
	if image == nil {
		return nil
	}
	return &Image{
		ID:       image.ID,
		W342:     image.W342,
		Original: image.Original,
	}
}

func mapTVShowShort(item themoviedb.TVShowShort) *TVShowShort {
	return &TVShowShort{
		ID:           item.ID,
		Name:         item.Name,
		OriginalName: item.OriginalName,
		Overview:     item.Overview,
		Poster:       mapImage(item.Poster),
		FirstAirDate: item.FirstAirDate,
		VoteAverage:  item.VoteAverage,
		VoteCount:    item.VoteCount,
		Popularity:   item.Popularity,
	}
}

func mapTVShowShorts(items []themoviedb.TVShowShort) []TVShowShort {
	return lo.Map(items, func(item themoviedb.TVShowShort, index int) TVShowShort {
		return *mapTVShowShort(item)
	})
}

func mapSeason(season themoviedb.Season) *Season {
	return &Season{
		ID:           season.ID,
		AirDate:      season.AirDate,
		EpisodeCount: season.EpisodeCount,
		Name:         season.Name,
		Overview:     season.Overview,
		Poster:       mapImage(season.Poster),
		SeasonNumber: season.SeasonNumber,
		VoteAverage:  season.VoteAverage,
	}
}

func mapTVShow(response *themoviedb.TVShow) *TVShow {
	return &TVShow{
		TVShowShort:      *mapTVShowShort(response.TVShowShort),
		Backdrop:         mapImage(response.Backdrop),
		Genres:           response.Genres,
		LastAirDate:      response.LastAirDate,
		NextEpisodeToAir: response.NextEpisodeToAir,
		NumberOfEpisodes: response.NumberOfEpisodes,
		NumberOfSeasons:  response.NumberOfSeasons,
		OriginCountry:    response.OriginCountry,
		Status:           response.Status,
		Tagline:          response.Tagline,
		Type:             response.Type,
		Seasons: lo.Map(response.Seasons, func(item themoviedb.Season, index int) Season {
			return *mapSeason(item)
		}),
	}
}

func mapEpisodes(response []themoviedb.Episode) []Episode {
	return lo.Map(response, func(item themoviedb.Episode, index int) Episode {
		return Episode{
			ID:            item.ID,
			AirDate:       item.AirDate,
			EpisodeNumber: item.EpisodeNumber,
			EpisodeType:   item.EpisodeType,
			Name:          item.Name,
			Overview:      item.Overview,
			Runtime:       item.Runtime,
			Still:         mapImage(item.Still),
			VoteAverage:   item.VoteAverage,
			VoteCount:     item.VoteCount,
		}
	})
}
