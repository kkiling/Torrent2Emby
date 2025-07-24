package mkvmerge

import "github.com/kkiling/torrent2emby/internal/log"

type Service struct {
	logger log.Logger
}

func NewService(logger log.Logger) *Service {
	return &Service{
		logger: logger.Named("mkvmerge"),
	}
}
