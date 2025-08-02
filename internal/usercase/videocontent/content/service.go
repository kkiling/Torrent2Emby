package content

import (
	"context"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type Service struct {
	storage Storage
}

func NewService(
	storage Storage,
) *Service {
	return &Service{
		storage: storage,
	}
}

// CreateVideoContent создание файловой раздачи
func (s *Service) CreateVideoContent(ctx context.Context, params CreateVideoContentParams) (*VideoContent, error) {
	// Добавить проверку на наличие уже раздачи и выкидывать ее
	panic("implement me")
	return nil, ucerr.NotFound
}

func (s *Service) GetVideoContent(ctx context.Context, contentID videocontent.ContentID) ([]VideoContent, error) {
	panic("implement me")
}

func (s *Service) GetTVShowDeliveryState(ctx context.Context, contentID videocontent.ContentID) (*tvshowdeliverystate.State, error) {
	panic("implement me")
	return nil, ucerr.NotFound
}

// Запускаем по cron
func (s *Service) updateDeliveryStatus(ctx context.Context, contentID videocontent.ContentID) {

	//Трекаем обновление статуса в процессе доставки in_progress до доставлено delivered на основе стейта
	//	in_progress -> delivered
	//Трекаем по аналогии
	//	updating -> delivered
	//Трекаем по анлогии на основании стейта
	//	deleting -> deleted

}
