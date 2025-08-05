package handler

import (
	"fmt"

	"github.com/kkiling/goplatform/server"

	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
)

// HandleError обработчик ошибок
func HandleError(err error, description any) error {
	newErr := fmt.Errorf("%s: %w", description, err)

	info := desc.ErrorInfo{
		Description: "Unhandled error",
	}

	return server.ErrInternal(newErr, &info)
}
