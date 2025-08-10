package handler

import (
	"errors"
	"fmt"

	"github.com/kkiling/goplatform/server"
	"github.com/kkiling/statemachine"

	"github.com/kkiling/torrent2emby/internal/usercase/err"
	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
)

// HandleError обработчик ошибок
func HandleError(err error, description any) error {
	newErr := fmt.Errorf("%s: %w", description, err)

	switch {
	case errors.Is(err, ucerr.NotFound):
		return server.ErrNotFound(newErr)
	case errors.Is(err, ucerr.InvalidArgument):
		return server.ErrInvalidArgument(newErr)
	case errors.Is(err, ucerr.AlreadyExists):
		return server.ErrAlreadyExists(newErr)
	case errors.Is(err, statemachine.ErrOptionsIsUndefined):
		return server.ErrInvalidArgument(newErr)
	case errors.Is(err, statemachine.ErrInTerminalStatus):
		return server.ErrAlreadyExists(newErr)
	case errors.Is(err, statemachine.ErrAlreadyExists):
		return server.ErrAlreadyExists(newErr)
	case errors.Is(err, statemachine.ErrNotFound):
		return server.ErrNotFound(newErr)

	}

	info := desc.ErrorInfo{
		Description: "Unhandled error",
	}

	return server.ErrInternal(newErr, &info)
}
