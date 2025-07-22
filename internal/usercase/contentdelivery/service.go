package contentdelivery

type Config struct {
	// BasePath Базовый путь от которого расположены все файлы торрента или медиа сервера
	// Например скачанные сериалы лежат по пути BasePath + TVShowTorrentSavePath
	BasePath string // "/nfs"
	// TVShowTorrentSavePath путь сохранения сериалов относительно торрент клиента
	TVShowTorrentSavePath string
	// TvShowMediaSaveTvShowsPath путь сохранения сериалов относительно медиа сервера
	TvShowMediaSaveTvShowsPath string
	// Группа от именни которой будут проводиться манипуляции с файлами и каталогами (если не указать, то будет исопользована дефолтная группа пользователя)
	UserGroup string
}

type Service struct {
	config        Config
	tvShowLibrary TVShowLibrary
	torrentSite   TorrentSite
	torrentClient TorrentClient
	prepareTVShow PrepareTVShow
	mkvMerge      MkvMerge
}

func NewService(
	config Config,
	tvShowLibrary TVShowLibrary,
	torrentSite TorrentSite,
	torrentClient TorrentClient,
	prepareTVShow PrepareTVShow,
	mkvMerge MkvMerge,
) *Service {
	return &Service{
		config:        config,
		tvShowLibrary: tvShowLibrary,
		torrentSite:   torrentSite,
		torrentClient: torrentClient,
		prepareTVShow: prepareTVShow,
		mkvMerge:      mkvMerge,
	}
}
