package videocontent

import (
	"context"
	"fmt"
)

type GetMagnetLinkParams struct {
	Href string
}

// getMagnetLink получение магнет ссылки на основе выбора раздачи пользователем
func (s *Service) getMagnetLink(ctx context.Context, params GetMagnetLinkParams) (*MagnetInfo, error) {
	// Получение магнет ссылки
	magnetInfo, err := s.torrentSite.GetMagnetLink(params.Href)
	if err != nil {
		return nil, fmt.Errorf("torrentSite.GetMagnetLink: %w", err)
	}

	return &MagnetInfo{
		Magnet: magnetInfo.Magnet,
		Hash:   magnetInfo.Hash,
	}, nil
}
