package delivery

import (
	"context"
	"fmt"

	getFolderSize "github.com/markthree/go-get-folder-size/src"
)

// GetCatalogSize получение размера каталога на диске в байтах
func (s *Service) GetCatalogSize(_ context.Context, catalogPath string) (uint64, error) {
	size, err := getFolderSize.Invoke(catalogPath)
	if err != nil {
		return 0, fmt.Errorf("could not get size: %w", err)
	}
	return uint64(size), nil
}
