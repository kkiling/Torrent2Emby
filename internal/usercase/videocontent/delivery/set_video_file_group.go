package delivery

import (
	"context"
	"fmt"
)

// SetVideoFileGroup установка группы указаным файлам
func (s *Service) SetVideoFileGroup(_ context.Context, files []string) error {
	if s.config.UserGroup == "" {
		return nil
	}

	for _, file := range files {
		if err := setGroup(file, s.config.UserGroup); err != nil {
			return fmt.Errorf("setting video file group: %w", err)
		}
	}

	return nil
}
