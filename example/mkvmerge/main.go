package main

import (
	"context"
	"errors"
	"fmt"
	mkvmerge "github.com/kkiling/torrent2emby/internal/adapter/mkvmerge"
	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/usercase/mkvmergepipeline"
	sqlite "github.com/kkiling/torrent2emby/internal/usercase/mkvmergepipeline/storage/sqlite"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogger(log.DebugLevel)

	store, err := sqlite.NewStorage(sqlite.Config{DSN: "/home/kiling/projects/torrent2emby/torrent2emby.db"}, logger)
	if err != nil {
		logger.Fatal(err)
	}

	mkv := mkvmerge.NewService(logger)
	service := mkvmergepipeline.NewService(mkvmergepipeline.Config{}, mkv, store, logger)

	idempotencyKey := "test-1"
	result, err := service.AddToMerge(ctx, idempotencyKey, mkvmerge.MergeParams{
		VideoInputFile:  "/nfs/downloads/Vinland Saga S2/Vinland Saga S2 [04].avi",
		VideoOutputFile: "/home/kiling/Desktop/sagatest/4.avi",
		AudioTracks:     nil,
		SubtitleTracks: []mkvmerge.Track{
			{
				Path:     "/nfs/downloads/Vinland Saga S2/Rus subs/Vinland Saga S2 [04].ass",
				Language: "rus",
				Name:     "Rus subs",
				Default:  true,
			},
		},
	})
	if err != nil {
		if errors.Is(err, mkvmergepipeline.ErrAlreadyExists) {
			fmt.Println("already exist")
		} else {
			logger.Fatal(err)
		}
	}

	fmt.Println(result.ID)

	err = service.StartMergePipeline(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	find, err := service.GetMergeResult(ctx, result.ID)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(find)
}
