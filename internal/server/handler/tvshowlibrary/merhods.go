package tvshowlibrary

import (
	"context"

	"github.com/kkiling/torrent2emby/internal/server/handler"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
)

func (h *Handler) SearchTVShow(ctx context.Context, request *desc.SearchTVShowRequest) (*desc.SearchTVShowResponse, error) {
	result, err := h.tvShowLibrary.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
		Query: request.Query,
	})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.SearchTVShow")
	}

	return &desc.SearchTVShowResponse{
		Items: mapTvShowShorts(result.Items),
	}, nil
}

func (h *Handler) GetTVShowInfo(ctx context.Context, request *desc.GetTVShowInfoRequest) (*desc.GetTVShowInfoResponse, error) {
	result, err := h.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: request.TvShowId,
	})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.GetTVShowInfo")
	}

	return &desc.GetTVShowInfoResponse{Result: mapTvShow(result.Result)}, nil
}

func (h *Handler) GetSeasonEpisodes(ctx context.Context, request *desc.GetSeasonEpisodesRequest) (*desc.GetSeasonEpisodesResponse, error) {
	result, err := h.tvShowLibrary.GetSeasonEpisodes(ctx, tvshowlibrary.GetSeasonEpisodesParams{
		TVShowID:     request.TvShowId,
		SeasonNumber: uint8(request.SeasonNumber),
	})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.GetSeasonEpisodes")
	}

	return &desc.GetSeasonEpisodesResponse{
		Items: mapEpisodes(result.Items),
	}, nil
}

func (h *Handler) GetTVShowsFromLibrary(ctx context.Context, request *desc.GetTVShowsFromLibraryRequest) (*desc.GetTVShowsFromLibraryResponse, error) {
	result, err := h.tvShowLibrary.GetTVShowsFromLibrary(ctx, tvshowlibrary.GetTVShowsFromLibraryParams{})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.GetSeasonEpisodes")
	}

	return &desc.GetTVShowsFromLibraryResponse{
		Items: mapTvShowShorts(result.Items),
	}, nil
}
