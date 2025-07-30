package contentdelivery

import (
	"context"
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/emby"
	"path/filepath"
)

type SetMediaMetaDataParams struct {
	SeasonPath   string
	TheMovieDBID uint64
}

// SetMediaMetaData установка методанных
func (s *Service) SetMediaMetaData(ctx context.Context, params SetMediaMetaDataParams) error {
	seasonPath, err := filepath.Rel(s.config.BasePath, params.SeasonPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	info, err := s.embyApi.GetCatalogInfo("/" + seasonPath)
	if err != nil {
		return fmt.Errorf("embyApi.GetCatalogInfo: %w", err)
	}
	if info == nil {
		return fmt.Errorf("catalogInfo: info is nil")
	}

	if !info.IsFolder {
		return fmt.Errorf("catalogInfo: info is not folder")
	}
	if info.Type != emby.SeriesTypeCatalog {
		return fmt.Errorf("catalogInfo: type is not series")
	}

	err = s.embyApi.RemoteSearchApply(info.ID, params.TheMovieDBID)
	if err != nil {
		return fmt.Errorf("embyApi.RemoteSearchApply: %w", err)
	}

	return nil
}
