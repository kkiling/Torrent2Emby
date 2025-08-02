package delivery

import (
	"context"
	"fmt"
)

type SearchTorrentParams struct {
	SearchQuery string
}

// SearchTorrent Делаем запрос к торрент сайту, получаем список раздач
func (s *Service) SearchTorrent(_ context.Context, params SearchTorrentParams) (*TorrentSearchResult, error) {
	searchResult, err := s.torrentSite.SearchTorrents(params.SearchQuery)
	if err != nil {
		return nil, fmt.Errorf("torrentSite.SearchTorrents: %w", err)
	}

	result := TorrentSearchResult{}
	for _, item := range searchResult.Results {
		result.Result = append(result.Result, TorrentSearch{
			Title:     item.Title,
			Href:      item.Href,
			Size:      item.Size,
			Seeds:     item.Seeds,
			Leeches:   item.Leeches,
			Downloads: item.Downloads,
			AddedDate: item.AddedDate,
		})
	}

	return &result, nil
}
