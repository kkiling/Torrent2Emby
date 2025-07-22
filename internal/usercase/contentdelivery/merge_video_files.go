package contentdelivery

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
)

type MergeVideoFilesParams struct {
	// Hash хеш торрента
	Hash           string
	ContentPath    string
	ContentMatches []ContentMatches
	// Курсор сколько обработано файлов, что бы вернуться к последнему
	ProcessedFiles int
}

func mapMkvMergeParams(content ContentMatches, contentPath string) mkvmerge.MergeParams {
	episodeName := fmt.Sprintf("S03%dE03%d %s", content.ContentInfo.SeasonNumber, content.ContentInfo.EpisodeNumber, content.ContentInfo.Name)

	mergeParams := mkvmerge.MergeParams{
		VideoInputFile: content.Video.File.FullPath,
		// Формирование исходного имени файла серии
		VideoOutputFile: filepath.Join(contentPath, episodeName) + content.Video.File.Extension,
		AudioTracks: lo.Map(content.AudioFiles, func(item Track, index int) mkvmerge.Track {
			return mkvmerge.Track{
				Path:     item.File.FullPath,
				Language: item.Language,
				Name:     item.Name,
				Default:  index == 0,
			}
		}),
		SubtitleTracks: lo.Map(content.Subtitles, func(item Track, index int) mkvmerge.Track {
			return mkvmerge.Track{
				Path:     item.File.FullPath,
				Language: item.Language,
				Name:     item.Name,
				Default:  false,
			}
		}),
	}
	return mergeParams
}

// MergeVideoFiles запуск обработки видеофайлов
func (s *Service) MergeVideoFiles(_ context.Context, params MergeVideoFilesParams) (MergeVideoStatus, error) {
	for index, content := range params.ContentMatches {
		if index < params.ProcessedFiles {
			continue
		}
		mergeParams := mapMkvMergeParams(content, params.ContentPath)
		err := s.mkvMerge.Merge(mergeParams)
		if err != nil {
			return MergeVideoStatus{
					ProcessedFiles: params.ProcessedFiles, // Возвращаем старое значение
					AllFiles:       len(params.ContentMatches),
					IsComplete:     false,
				},
				fmt.Errorf("mkvMerge.Merge: %w", err)
		}

		info, err := s.mkvMerge.GetMediaInfo(mergeParams.VideoOutputFile)
		if err != nil {
			return MergeVideoStatus{
					ProcessedFiles: params.ProcessedFiles, // Возвращаем старое значение
					AllFiles:       len(params.ContentMatches),
					IsComplete:     false,
				},
				fmt.Errorf("mkvMerge.GetMediaInfo: %w", err)
		}

		// TODO: Валидация мержинга файлов
		fmt.Println(info)

		// Выходим после каждого обработанного файла, что бы сохранить состояние
		return MergeVideoStatus{
			ProcessedFiles: index + 1,
			AllFiles:       len(params.ContentMatches),
			IsComplete:     index == len(params.ContentMatches)-1,
		}, nil
	}

	return MergeVideoStatus{}, fmt.Errorf("unknow merge statemachine")
}
