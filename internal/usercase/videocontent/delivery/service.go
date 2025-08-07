package delivery

type Config struct {
	// BasePath Базовый путь от которого расположены все файлы торрента или медиа сервера
	// Например скачанные сериалы лежат по пути BasePath + TVShowTorrentSavePath
	BasePath string // "/nfs"
	// TVShowTorrentSavePath путь сохранения сериалов относительно торрент клиента
	TVShowTorrentSavePath string
	// TVShowMediaSaveTvShowsPath путь сохранения сериалов относительно медиа сервера
	TVShowMediaSaveTvShowsPath string
	// Группа от именни которой будут проводиться манипуляции с файлами и каталогами (если не указать, то будет исопользована дефолтная группа пользователя)
	UserGroup string
}

type Service struct {
	config        Config
	tvShowLibrary TVShowLibrary
	torrentSite   TorrentSite
	torrentClient TorrentClient
	embyApi       EmbyApi
	prepareTVShow PrepareTVShow
	mkvMerge      MkvMergePipeline
}

func NewService(
	config Config,
	tvShowLibrary TVShowLibrary,
	torrentSite TorrentSite,
	torrentClient TorrentClient,
	embyApi EmbyApi,
	prepareTVShow PrepareTVShow,
	mkvMerge MkvMergePipeline,
) *Service {
	return &Service{
		config:        config,
		tvShowLibrary: tvShowLibrary,
		torrentSite:   torrentSite,
		torrentClient: torrentClient,
		embyApi:       embyApi,
		prepareTVShow: prepareTVShow,
		mkvMerge:      mkvMerge,
	}
}
