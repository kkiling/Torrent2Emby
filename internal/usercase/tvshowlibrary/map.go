package tvshowlibrary

import (
	"github.com/kkiling/torrent2emby/internal/adapter/themoviedb"
	"github.com/samber/lo"
)

func mapImage(image *themoviedb.Image) *Image {
	return &Image{
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
		PosterPath:   mapImage(item.PosterPath),
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
		PosterPath:   mapImage(season.PosterPath),
		SeasonNumber: season.SeasonNumber,
		VoteAverage:  season.VoteAverage,
	}
}

func mapTVShow(response *themoviedb.TVShow) *TVShow {
	return &TVShow{
		TVShowShort:      *mapTVShowShort(response.TVShowShort),
		BackdropPath:     mapImage(response.BackdropPath),
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
			StillPath:     mapImage(item.StillPath),
			VoteAverage:   item.VoteAverage,
			VoteCount:     item.VoteCount,
		}
	})
}
