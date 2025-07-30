package tvshowdeliverystate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/statemachine"
	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowdelivery"
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
	// Логика создания задачи
	data := ContentDeliveryData{}

	return CreateState{
		FirstStep: GenerateSearchQuery,
		Data:      data,
		MetaData: ContentDeliveryMetadata{
			TVShowID: options.TVShowID,
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
					res, err := r.contentDelivery.GenerateSearchQuery(ctx, tvshowdelivery.GenerateSearchQueryParams{
						TVShowID: stepContext.State.MetaData.TVShowID,
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
					res, err := r.contentDelivery.SearchTorrent(ctx, tvshowdelivery.SearchTorrentParams{
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
						contains := lo.ContainsBy(data.TorrentSearch.Result, func(item tvshowdelivery.TorrentSearch) bool {
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
					res, err := r.contentDelivery.GetMagnetLink(ctx, tvshowdelivery.GetMagnetLinkParams{
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
					err := r.contentDelivery.AddTorrentToTorrentClient(ctx, tvshowdelivery.AddTorrentParams{
						TVShowID: stepContext.State.MetaData.TVShowID,
						Magnet:   data.MagnetInfo.Magnet,
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
					res, err := r.contentDelivery.PrepareFileMatches(ctx, tvshowdelivery.PreparingFileMatchesParams{
						Hash:     data.MagnetInfo.Hash,
						TVShowID: stepContext.State.MetaData.TVShowID,
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
					res, err := r.contentDelivery.WaitingTorrentDownloadComplete(ctx, tvshowdelivery.WaitingTorrentDownloadCompleteParams{
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
					res, err := r.contentDelivery.CreateContentCatalogs(ctx, tvshowdelivery.CreateContentCatalogsParams{
						TVShowID: stepContext.State.MetaData.TVShowID,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("CreateContentCatalogs: %w", err))
					}
					data := stepContext.State.Data
					data.CatalogsInfo = res
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
						return stepContext.Next(StartMergeVideoFiles)
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
			StartMergeVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					data := stepContext.State.Data

					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					mergeIDs, err := r.contentDelivery.StartMergeVideo(ctx, tvshowdelivery.MergeVideoParams{
						ContentPath:    data.CatalogsInfo.TvShowSeasonPath,
						IdempotencyKey: stepContext.State.ID.String(),
						ContentMatches: data.ContentMatches,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("StartMergeVideoFiles: %w", err))
					}
					data.MergeVideoFiles = mergeIDs
					return stepContext.Next(WaitingMergeVideoFiles).WithData(data)
				},
			},
			WaitingMergeVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					mergeIDs := lo.Map(data.MergeVideoFiles, func(item tvshowdelivery.MergeVideoFile, _ int) uuid.UUID {
						return item.MergeID
					})
					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					status, err := r.contentDelivery.GetMergeVideoStatus(ctx, mergeIDs)
					if err != nil {
						return stepContext.Error(fmt.Errorf("WaitingMergeVideoFiles: %w", err))
					}

					data.MergeVideoStatus = status
					if status.IsComplete {
						if len(status.Errors) == 0 {
							return stepContext.Next(SetVideoFileGroup).WithData(data)
						}
						return stepContext.Error(fmt.Errorf("merge videos contains errors")).WithData(data)
					}
					return stepContext.Empty().WithData(data)
				},
			},
			SetVideoFileGroup: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data
					files := lo.Map(data.MergeVideoFiles, func(item tvshowdelivery.MergeVideoFile, _ int) string {
						return item.VideoOutputFile
					})

					// Установка группы файлам
					err := r.contentDelivery.SetVideoFileGroup(ctx, files)
					if err != nil {
						return stepContext.Error(fmt.Errorf("SetVideoFileGroup: %w", err))
					}

					return stepContext.Next(SetMediaMetaData)
				},
			},
			SetMediaMetaData: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data
					// Установка группы файлам
					err := r.contentDelivery.SetMediaMetaData(ctx, tvshowdelivery.SetMediaMetaDataParams{
						SeasonPath: data.CatalogsInfo.TvShowPath,
						TVShowID:   stepContext.State.MetaData.TVShowID,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("SetVideoFileGroup: %w", err))
					}

					return stepContext.Complete().WithData(data)
				},
			},
		},
	}
}
