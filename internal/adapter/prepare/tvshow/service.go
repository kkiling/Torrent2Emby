package tvshow

import (
	"fmt"
	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	"path/filepath"
	"strings"
)

var (
	videoExtensions     = []string{".mp4", ".mkv", ".avi", ".mp3", ".flac", ".m4a", ".ogg", ".wav"}
	audioExtensions     = []string{".mka"}
	subtitlesExtensions = []string{".ass"}
)

type MediaInfo interface {
	GetMediaInfo(filePath string) (*mkvmerge.MediaInfo, error)
}

type Service struct {
	mediaInfo MediaInfo
}

func NewService(mediaInfo MediaInfo) *Service {
	return &Service{
		mediaInfo: mediaInfo,
	}
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

		// Проверяем расширение файла
		if !extMap[file.Extension] {
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

func (s *Service) processMetaFiles(torrentFiles []TorrentFile, extensions []string) (map[string][]PrepareTrack, error) {
	prepareVideoFiles, err := processFiles(torrentFiles, extensions)
	if err != nil {
		return nil, fmt.Errorf("processFiles: %w", err)
	}
	// Группируем по озвучке
	result := make(map[string][]PrepareTrack)
	for _, file := range prepareVideoFiles {
		// Исходим из того что озвучка/субтитры лежит в каком то каталоге
		// Название этого каталога и берем за название озвучки/субтитры
		splitRelativePath := splitPath(file.RelativePath)
		if len(splitRelativePath) < 2 {
			continue
		}

		// Пробуем достать информацию из файла
		info, err := s.mediaInfo.GetMediaInfo(file.FullPath)
		if err != nil {
			return nil, fmt.Errorf("mediaInfo.GetMediaInfo: %w", err)
		}

		// Берем название из каталога
		name := splitRelativePath[len(splitRelativePath)-2]
		language := ""
		if len(info.AudioTracks) == 1 && info.AudioTracks[0].TrackName != "" {
			// Пробуем достать из инфы аудиодорожки
			name = info.AudioTracks[0].TrackName
			language = info.AudioTracks[0].Language
		}

		if len(info.Subtitles) == 1 && info.Subtitles[0].TrackName != "" {
			// Пробуем достать из инфы субтитров
			name = info.Subtitles[0].TrackName
			language = info.AudioTracks[0].Language
		}

		result[name] = append(result[name], PrepareTrack{
			Name:     name,
			Language: language,
			File:     file,
		})
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
	audioFilesMap, err := s.processMetaFiles(params.TorrentFiles, audioExtensions)
	if err != nil {
		return nil, fmt.Errorf("processFiles audio files: %w", err)
	}

	// Получаем субтитры
	subtitlesFilesMap, err := s.processMetaFiles(params.TorrentFiles, subtitlesExtensions)
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

		for _, audioFiles := range audioFilesMap {
			if index >= len(audioFiles) {
				continue
			}
			prepareEpisode.AudioFiles = append(prepareEpisode.AudioFiles, audioFiles[index])
		}

		for _, subtitleFiles := range subtitlesFilesMap {
			if index >= len(subtitleFiles) {
				continue
			}
			prepareEpisode.Subtitles = append(prepareEpisode.Subtitles, subtitleFiles[index])
		}

		result.Episodes = append(result.Episodes, prepareEpisode)
	}

	return &result, nil
}
