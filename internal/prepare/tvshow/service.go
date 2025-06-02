package tvshow

import (
	"fmt"
	"path/filepath"
	"strings"
)

var (
	videoExtensions     = []string{".mp4", ".mkv", ".avi", ".mp3", ".flac", ".m4a", ".ogg", ".wav"}
	audioExtensions     = []string{".mka"}
	subtitlesExtensions = []string{".ass"}
)

type Service struct {
	// embyMediaPath каталог с которым работает emby сервер, куда нужно копировать файлы
	embyMediaPath string
}

func NewService() *Service {
	return &Service{}
}

// splitPath разбивает путь в Linux на отдельные компоненты.
// Пример:
//
//	"/home/user/docs/file.txt" -> ["home", "user", "docs", "file.txt"]
//	"relative/path/" -> ["relative", "path"]
//	"/" -> []
func splitPath(path string) []string {
	// Нормализуем путь (убираем дублирующиеся слеши, обработка . и ..)
	cleanPath := filepath.Clean(path)

	// Разбиваем на компоненты
	parts := strings.Split(cleanPath, "/")

	// Удаляем пустые элементы (могут появиться из-за концевого слеша)
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

// processFile на основе торрент файлов получаем спиос кафйлов нужных расшерений
func processFiles(
	files []TorrentFile,
	allowedExtensions []string,
) ([]TorrentFile, error) {

	// Создаем map для быстрой проверки расширений
	extMap := make(map[string]bool)
	for _, ext := range allowedExtensions {
		extMap[strings.ToLower(ext)] = true
	}

	var result []TorrentFile
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.RelativePath))

		// Проверяем расширение файла
		if !extMap[ext] {
			continue
		}

		// Добавляем в результат
		result = append(result, file)
	}

	return result, nil
}

func processVideoFiles(torrentFiles []TorrentFile) ([]TorrentFile, error) {
	prepareVideoFiles, err := processFiles(torrentFiles, videoExtensions)
	if err != nil {
		return nil, fmt.Errorf("processFiles: %w", err)
	}

	result := make([]TorrentFile, 0)
	for _, file := range prepareVideoFiles {
		// Исходим из того что файлы видео файлов серий всегда лежат в корне
		splitRelativePath := splitPath(file.RelativePath)
		if len(splitRelativePath) > 1 {
			continue
		}

		result = append(result, file)
	}

	return result, nil
}

func processMetaFiles(torrentFiles []TorrentFile, extensions []string) (map[string][]TorrentFile, error) {
	prepareVideoFiles, err := processFiles(torrentFiles, extensions)
	if err != nil {
		return nil, fmt.Errorf("processFiles: %w", err)
	}
	// Группируем по озвучке
	result := make(map[string][]TorrentFile)
	for _, file := range prepareVideoFiles {
		// Исходим из того что озвучка/субтитры лежит в каком то каталоге
		// Название этого каталога и берем за название озвучки/субтитры
		splitRelativePath := splitPath(file.RelativePath)
		if len(splitRelativePath) < 2 {
			continue
		}
		// Берем название
		name := splitRelativePath[len(splitRelativePath)-2]
		result[name] = append(result[name], file)
	}

	return result, nil
}

func (s *Service) PrepareTvShowSeason(params *PrepareTvShowPrams) (*PrepareTVShowSeason, error) {
	result := PrepareTVShowSeason{}

	// Получаем список видео файлов эпизодов
	prepareVideoFiles, err := processVideoFiles(params.TorrentFiles)
	if err != nil {
		return nil, fmt.Errorf("processFiles: %w", err)
	}

	// Получаем аудиодорожки
	audioFilesMap, err := processMetaFiles(params.TorrentFiles, audioExtensions)
	if err != nil {
		return nil, fmt.Errorf("processFiles audio files: %w", err)
	}

	// Получаем субтитры
	subtitlesFilesMap, err := processMetaFiles(params.TorrentFiles, subtitlesExtensions)
	if err != nil {
		return nil, fmt.Errorf("processFiles subtitles files: %w", err)
	}

	for index, episode := range params.Episodes {

		prepareEpisode := PrepareEpisode{
			Episode: episode,
		}

		// Пока сопостовляем видео файл с серией просто по порядку
		if index < len(prepareVideoFiles) {
			prepareEpisode.VideoFile = &PrepareVideo{
				File: prepareVideoFiles[index],
			}
		}

		for audioName, audioFiles := range audioFilesMap {
			if index < len(audioFiles) {
				prepareEpisode.AudioFiles = append(prepareEpisode.AudioFiles, PrepareAudio{
					Name: audioName,
					File: audioFiles[index],
				})
			}
		}

		for subtitleName, subtitleFiles := range subtitlesFilesMap {
			if index < len(subtitleFiles) {
				prepareEpisode.Subtitles = append(prepareEpisode.Subtitles, PrepareSubtitles{
					Name: subtitleName,
					File: subtitleFiles[index],
				})
			}
		}

		result.Episodes = append(result.Episodes, prepareEpisode)
	}

	return &result, nil
}
