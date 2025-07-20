package runner

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kkiling/torrent2emby/internal/statemachine"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/videodelivery"
)

type Runner struct {
	contentDelivery ContentDelivery
}

func NewTaskRunner(contentDelivery ContentDelivery) *Runner {
	return &Runner{
		contentDelivery: contentDelivery,
	}
}

func (r *Runner) Create(_ context.Context, options CreateOptions) (CreateState, error) {
	if options.MediaID.MovieID == nil && options.MediaID.TVShow == nil {
		return CreateState{}, fmt.Errorf("movieID or TVShow is required: %w", ucerr.InvalidArgument)
	}

	// Логика создания задачи
	data := ContentDeliveryData{}

	return CreateState{
		FirstStep: GenerateSearchQuery,
		Data:      data,
		MetaData: ContentDeliveryMetadata{
			MediaID: options.MediaID,
		},
	}, nil
}

func (r *Runner) Type() Type {
	return ContentDeliveryType
}

func (r *Runner) StepRegistration(_ statemachine.StepRegistrationParams) StepRegistration {
	return StepRegistration{
		Steps: map[StepDelivery]Step{
			GenerateSearchQuery: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Генерация запроса
					data := stepContext.State.Data
					res, err := r.contentDelivery.GenerateSearchQuery(ctx, videodelivery.GenerateSearchQueryParams{
						MediaID: stepContext.State.MetaData.MediaID,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("GenerateSearchQuery: %w", err))
					}
					data.SearchQuery = &res
					return stepContext.Next(SearchTorrents).WithData(data)
				},
			},
			SearchTorrents: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// ищем раздачи сезона сериала / фильма
					data := stepContext.State.Data
					res, err := r.contentDelivery.SearchTorrent(ctx, videodelivery.SearchTorrentParams{
						SearchQuery: *data.SearchQuery,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("SearchTorrent: %w", err))
					}
					data.TorrentSearch = res
					return stepContext.Next(WaitingUserChoseTorrent).WithData(data)
				},
			},
			WaitingUserChoseTorrent: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Ожидаем когда пользователь выберет раздачу
					// Или ожидаем что клиент изменит поисковый запрос, тогда прыгаем на SearchTorrents
					// Получение опций выполнения выпуска
					opts := ChoseTorrentOptions{}
					ok, err := stepContext.GetOptions(&opts)
					if err != nil {
						return stepContext.Error(err)
					}
					if !ok { // Пока не получили опцию, не идем дальше
						return stepContext.Empty()
					}
					if opts.Href == nil && opts.NewSearchQuery == nil {
						return stepContext.Error(fmt.Errorf("either Href or NewSearchQuery must be specified: %w", ucerr.InvalidArgument))
					} else if opts.Href != nil && opts.NewSearchQuery != nil {
						return stepContext.Error(fmt.Errorf("only one of Href or NewSearchQuery can be specified: %w", ucerr.InvalidArgument))
					}

					data := stepContext.State.Data
					if opts.NewSearchQuery != nil {
						// Снова производим поиск по раздачам
						data.SearchQuery = opts.NewSearchQuery
						return stepContext.Next(SearchTorrents).WithData(data)
					}
					if opts.Href != nil {
						// Пользователь выбрал раздачу для скачивания
						data.SelectTorrentHref = opts.Href
						return stepContext.Next(GetMagnetLink).WithData(data)
					}
					return stepContext.Error(fmt.Errorf("unknow state: %w", ucerr.InvalidArgument))
				},
				OptionsType: reflect.TypeOf(ChoseTorrentOptions{}),
			},
			GetMagnetLink: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение магнет ссылки
					data := stepContext.State.Data
					res, err := r.contentDelivery.GetMagnetLink(ctx, videodelivery.GetMagnetLinkParams{
						Href: *data.SelectTorrentHref,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("GetMagnetLink: %w", err))
					}
					data.MagnetInfo = res
					return stepContext.Next(AddTorrentToTorrentClient).WithData(data)
				},
			},
			AddTorrentToTorrentClient: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Добавление раздачи для скачивания торрент клиентом
					data := stepContext.State.Data
					res, err := r.contentDelivery.GetMagnetLink(ctx, videodelivery.GetMagnetLinkParams{
						Href: *data.SelectTorrentHref,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("GetMagnetLink: %w", err))
					}
					data.MagnetInfo = res
					return stepContext.Next(PrepareFileMatches).WithData(data)
				},
			},
			PrepareFileMatches: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение информации о файлах раздачи
					data := stepContext.State.Data
					res, err := r.contentDelivery.PrepareFileMatches(ctx, videodelivery.PreparingFileMatchesParams{
						Hash:    data.MagnetInfo.Hash,
						MediaID: stepContext.State.MetaData.MediaID,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("PrepareFileMatches: %w", err))
					}
					if len(res) == 0 {
						return stepContext.Empty()
					}
					data.ContentMatches = res
					return stepContext.Next(WaitingChoseFileMatches).WithData(data)
				},
			},
		},
	}
}
