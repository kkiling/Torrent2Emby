package deliverystate

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"reflect"

	"github.com/kkiling/torrent2emby/internal/statemachine"
	"github.com/kkiling/torrent2emby/internal/usercase/contentdelivery"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
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
					res, err := r.contentDelivery.GenerateSearchQuery(ctx, contentdelivery.GenerateSearchQueryParams{
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
					res, err := r.contentDelivery.SearchTorrent(ctx, contentdelivery.SearchTorrentParams{
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
						// Проверяем что клиент выбрал href из списка
						contains := lo.ContainsBy(data.TorrentSearch.Result, func(item contentdelivery.TorrentSearch) bool {
							return item.Href == *opts.Href
						})
						if !contains {
							return stepContext.Error(fmt.Errorf("no such href: %w", ucerr.InvalidArgument))
						}

						data.SelectTorrentHref = opts.Href
						return stepContext.Next(GetMagnetLink).WithData(data)
					}
					return stepContext.Error(fmt.Errorf("unknow deliverystate: %w", ucerr.InvalidArgument))
				},
				OptionsType: reflect.TypeOf(ChoseTorrentOptions{}),
			},
			GetMagnetLink: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение магнет ссылки
					data := stepContext.State.Data
					res, err := r.contentDelivery.GetMagnetLink(ctx, contentdelivery.GetMagnetLinkParams{
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
					err := r.contentDelivery.AddTorrentToTorrentClient(ctx, contentdelivery.AddTorrentParams{
						MediaID: stepContext.State.MetaData.MediaID,
						Magnet:  data.MagnetInfo.Magnet,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("AddTorrentToTorrentClient: %w", err))
					}
					return stepContext.Next(PrepareFileMatches).WithData(data)
				},
			},
			PrepareFileMatches: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение информации о файлах раздачи
					data := stepContext.State.Data
					res, err := r.contentDelivery.PrepareFileMatches(ctx, contentdelivery.PreparingFileMatchesParams{
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
			WaitingChoseFileMatches: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// ожидание подтверждения пользователем соответствий выбора файлов
					opts := ChoseFileMatchesOptions{}
					ok, err := stepContext.GetOptions(&opts)
					if err != nil {
						return stepContext.Error(err)
					}
					if !ok { // Пока не получили опцию, не идем дальше
						return stepContext.Empty()
					}
					if !opts.Approve {
						return stepContext.Empty()
					}
					// TODO: выбор пользовтелем другого сопоставления

					return stepContext.Next(WaitingTorrentDownloadComplete)
				},
				OptionsType: reflect.TypeOf(ChoseFileMatchesOptions{}),
			},
			WaitingTorrentDownloadComplete: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Ожидание когда торрент докачается до конца
					data := stepContext.State.Data
					res, err := r.contentDelivery.WaitingTorrentDownloadComplete(ctx, contentdelivery.WaitingTorrentDownloadCompleteParams{
						Hash: data.MagnetInfo.Hash,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("PrepareFileMatches: %w", err))
					}
					data.TorrentDownloadStatus = res
					if res.IsComplete {
						return stepContext.Next(CreateVideoContentCatalogs).WithData(data)
					}
					return stepContext.Empty().WithData(data)
				},
			},
			CreateVideoContentCatalogs: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Формирование каталогов и иерархии файлов
					res, err := r.contentDelivery.CreateContentCatalogs(ctx, contentdelivery.CreateContentCatalogsParams{
						MediaID: stepContext.State.MetaData.MediaID,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("CreateContentCatalogs: %w", err))
					}
					data := stepContext.State.Data
					data.CatalogsInfo = &res
					return stepContext.Next(DeterminingNeedConvertFiles).WithData(data)
				},
			},
			DeterminingNeedConvertFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Определение необходимости конвертации файлов

					data := stepContext.State.Data
					needToMerge := false
					for _, m := range data.ContentMatches {
						// Если есть субтитры или аудиодорожки то нужно мержить
						if len(m.AudioFiles) > 0 || len(m.Subtitles) > 0 {
							needToMerge = true
							break
						}
					}
					if needToMerge {
						return stepContext.Next(MergeVideoFiles)
					}
					return stepContext.Next(CopyVideoFiles)
				},
			},
			CopyVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
					// TODO: реализовать
					return stepContext.Empty()
				},
			},
			MergeVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					data := stepContext.State.Data

					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					result, err := r.contentDelivery.MergeVideoFiles(ctx, contentdelivery.MergeVideoFilesParams{
						Hash:           data.MagnetInfo.Hash,
						ContentPath:    data.CatalogsInfo.CatalogPath,
						ContentMatches: data.ContentMatches,
						ProcessedFiles: func() int { // Стартуем с последнего
							if data.MergeVideoStatus != nil {
								return data.MergeVideoStatus.ProcessedFiles
							}
							return 0
						}(),
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("MergeVideoFiles: %w", err))
					}

					data.MergeVideoStatus = &result
					if result.IsComplete {
						// Переход на следующий шаг
						return stepContext.Next(SetMediaMetaData).WithData(data)
					}
					return stepContext.Empty().WithData(data)
				},
			},
		},
	}
}
