package tvshowlibrary

import (
	"context"

	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
)

func (h *Handler) SearchTVShow(ctx context.Context, request *desc.SearchTVShowRequest) (*desc.SearchTVShowResponse, error) {
	return &desc.SearchTVShowResponse{
		Items: []*desc.TVShowShort{
			{
				Id:   12,
				Name: request.Query,
			},
		},
	}, nil
}
