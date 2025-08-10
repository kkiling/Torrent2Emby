package tvshowdeliverystate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/kkiling/statemachine"
	"github.com/samber/lo"

	ucerr "github.com/kkiling/torrent2emby/internal/usercase/err"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/delivery"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners"
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
	data := TVShowDeliveryData{}

	return CreateState{
		FirstStep: GenerateSearchQuery,
		Data:      data,
		MetaData: runners.Metadata{
			ContentID: common.ContentID{
				TVShow: &options.TVShowID,
			},
		},
	}, nil
}

func (r *Runner) Type() runners.Type {
	return runners.TVShowDelivery
}

func (r *Runner) StepRegistration(_ statemachine.StepRegistrationParams) StepRegistration {
	return StepRegistration{
		Steps: map[StepDelivery]Step{
			GenerateSearchQuery: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Генерация запроса
					data := stepContext.State.Data
					res, err := r.contentDelivery.GenerateSearchQuery(ctx, delivery.GenerateSearchQueryParams{
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
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
					res, err := r.contentDelivery.SearchTorrent(ctx, delivery.SearchTorrentParams{
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
						contains := lo.ContainsBy(data.TorrentSearch.Result, func(item delivery.TorrentSearch) bool {
							return item.Href == *opts.Href
						})
						if !contains {
							return stepContext.Error(fmt.Errorf("no such href: %w", ucerr.InvalidArgument))
						}

						data.TorrentInfo = &delivery.TorrentInfo{
							Href: *opts.Href,
						}
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
					res, err := r.contentDelivery.GetMagnetLink(ctx, delivery.GetMagnetLinkParams{
						Href: data.TorrentInfo.Href,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("GetMagnetLink: %w", err))
					}
					data.TorrentInfo = res
					return stepContext.Next(AddTorrentToTorrentClient).WithData(data)
				},
			},
			AddTorrentToTorrentClient: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Добавление раздачи для скачивания торрент клиентом
					data := stepContext.State.Data
					err := r.contentDelivery.AddTorrentToTorrentClient(ctx, delivery.AddTorrentParams{
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
						Magnet:   data.TorrentInfo.Magnet,
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
					res, err := r.contentDelivery.PrepareFileMatches(ctx, delivery.PreparingFileMatchesParams{
						Hash:     data.TorrentInfo.Hash,
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
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
					res, err := r.contentDelivery.WaitingTorrentDownloadComplete(ctx, delivery.WaitingTorrentDownloadCompleteParams{
						Hash: data.TorrentInfo.Hash,
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
					data := stepContext.State.Data

					res, err := r.contentDelivery.CreateContentCatalogs(ctx, delivery.CreateContentCatalogsParams{
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("CreateContentCatalogs: %w", err))
					}

					data.TVShowCatalogInfo = &delivery.TVShowCatalog{
						TorrentPath:     data.TorrentDownloadStatus.ContentPath,
						MediaServerPath: *res,
					}
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
						data.TVShowCatalogInfo.IsCopyFilesInMediaServer = true
						return stepContext.Next(StartMergeVideoFiles).WithData(data)
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

					mergeVideoFiles, err := r.contentDelivery.StartMergeVideo(ctx, delivery.MergeVideoParams{
						ContentPath:    data.TVShowCatalogInfo.MediaServerPath.FullSeasonPath(),
						IdempotencyKey: stepContext.State.ID.String(),
						ContentMatches: data.ContentMatches,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("StartMergeVideoFiles: %w", err))
					}
					data.MergeVideoFiles = mergeVideoFiles
					return stepContext.Next(WaitingMergeVideoFiles).WithData(data)
				},
			},
			WaitingMergeVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					mergeIDs := lo.Map(data.MergeVideoFiles, func(item delivery.MergeVideoFile, _ int) uuid.UUID {
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
					files := lo.Map(data.MergeVideoFiles, func(item delivery.MergeVideoFile, _ int) string {
						return item.VideoOutputFile
					})

					// Установка группы файлам
					err := r.contentDelivery.SetVideoFileGroup(ctx, files)
					if err != nil {
						return stepContext.Error(fmt.Errorf("SetVideoFileGroup: %w", err))
					}

					return stepContext.Next(GetCatalogsSize)
				},
			},
			GetCatalogsSize: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					// Установка группы файлам
					torrentTVShowSeasonSize, err := r.contentDelivery.GetCatalogSize(ctx, data.TVShowCatalogInfo.TorrentPath)
					if err != nil {
						return stepContext.Error(fmt.Errorf("contentDelivery.GetCatalogSize: %w", err))
					}
					mediaServerTVShowSize, err := r.contentDelivery.GetCatalogSize(ctx, data.TVShowCatalogInfo.MediaServerPath.FullSeasonPath())
					if err != nil {
						return stepContext.Error(fmt.Errorf("contentDelivery.GetCatalogSize: %w", err))
					}

					data.TVShowCatalogInfo.TorrentSize = torrentTVShowSeasonSize
					data.TVShowCatalogInfo.MediaServerSize = mediaServerTVShowSize

					return stepContext.Next(SetMediaMetaData).WithData(data)
				},
			},
			SetMediaMetaData: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data
					// Установка группы файлам
					err := r.contentDelivery.SetMediaMetaData(ctx, delivery.SetMediaMetaDataParams{
						SeasonPath: data.TVShowCatalogInfo.MediaServerPath.FullSeasonPath(),
						TVShowID:   *stepContext.State.MetaData.ContentID.TVShow,
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
