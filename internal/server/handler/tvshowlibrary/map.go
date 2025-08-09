package tvshowlibrary

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
)

func mapImage(res *tvshowlibrary.Image) *desc.Image {
	if res == nil {
		return nil
	}
	return &desc.Image{
		Id:       res.ID,
		W342:     res.W342,
		Original: res.Original,
	}
}

func mapTvShowShort(res tvshowlibrary.TVShowShort) *desc.TVShowShort {
	return &desc.TVShowShort{
		Id:           res.ID,
		Name:         res.Name,
		OriginalName: res.OriginalName,
		Overview:     res.Overview,
		Poster:       mapImage(res.Poster),
		FirstAirDate: timestamppb.New(res.FirstAirDate),
		VoteAverage:  float32(res.VoteAverage),
		VoteCount:    uint32(res.VoteCount),
		Popularity:   float32(res.Popularity),
	}
}

func mapTvShowShorts(res []tvshowlibrary.TVShowShort) []*desc.TVShowShort {
	return lo.Map(res, func(item tvshowlibrary.TVShowShort, _ int) *desc.TVShowShort {
		return mapTvShowShort(item)
	})
}

func mapSeason(res tvshowlibrary.Season) *desc.Season {
	return &desc.Season{
		Id:           res.ID,
		AirDate:      timestamppb.New(res.AirDate),
		EpisodeCount: res.EpisodeCount,
		Name:         res.Name,
		Overview:     res.Overview,
		Poster:       mapImage(res.Poster),
		SeasonNumber: uint32(res.SeasonNumber),
		VoteAverage:  float32(res.VoteAverage),
	}
}

func mapTvShow(res *tvshowlibrary.TVShow) *desc.TVShow {
	return &desc.TVShow{
		Id:               res.ID,
		Name:             res.Name,
		OriginalName:     res.OriginalName,
		Overview:         res.Overview,
		Poster:           mapImage(res.Poster),
		FirstAirDate:     timestamppb.New(res.FirstAirDate),
		VoteAverage:      float32(res.VoteAverage),
		VoteCount:        res.VoteCount,
		Popularity:       float32(res.Popularity),
		Backdrop:         mapImage(res.Backdrop),
		Genres:           res.Genres,
		LastAirDate:      timestamppb.New(res.LastAirDate),
		NextEpisodeToAir: timestamppb.New(res.NextEpisodeToAir),
		NumberOfEpisodes: res.NumberOfEpisodes,
		NumberOfSeasons:  res.NumberOfSeasons,
		OriginCountry:    res.OriginCountry,
		Status:           res.Status,
		Tagline:          res.Tagline,
		Type:             res.Type,
		Seasons: lo.Map(res.Seasons, func(item tvshowlibrary.Season, _ int) *desc.Season {
			return mapSeason(item)
		}),
	}
}

func mapEpisode(res tvshowlibrary.Episode) *desc.Episode {
	return &desc.Episode{
		Id:            res.ID,
		AirDate:       timestamppb.New(res.AirDate),
		EpisodeNumber: uint32(res.EpisodeNumber),
		EpisodeType:   res.EpisodeType,
		Name:          res.Name,
		Overview:      res.Overview,
		Runtime:       uint32(res.Runtime),
		Still:         mapImage(res.Still),
		VoteAverage:   float32(res.VoteAverage),
		VoteCount:     uint32(res.VoteCount),
	}
}

func mapEpisodes(res []tvshowlibrary.Episode) []*desc.Episode {
	return lo.Map(res, func(item tvshowlibrary.Episode, _ int) *desc.Episode {
		return mapEpisode(item)
	})
}
