package mkvmerge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type MediaInfo struct {
	FileName    string
	Duration    float64 // в секундах
	VideoTracks []MediaInfoTrack
	AudioTracks []MediaInfoTrack
	Subtitles   []MediaInfoTrack
}

type MediaInfoTrack struct {
	Number            int
	DefaultDuration   int
	DefaultTrack      bool
	DisplayDimensions string
	EnabledTrack      bool
	ForcedTrack       bool
	Language          string
	TrackName         string
}

func (s *Merge) GetMediaInfo(filePath string) (*MediaInfo, error) {
	cmd := exec.Command("mkvmerge", "-J", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run mkvmerge: %v", err)
	}

	var result struct {
		FileName  string `json:"file_name"`
		Container struct {
			Properties struct {
				Duration int64 `json:"duration"`
			} `json:"properties"`
		} `json:"container"`
		Tracks []struct {
			Codec      string `json:"codec"`
			Id         int    `json:"id"`
			Properties struct {
				Number            int    `json:"number"`
				DefaultDuration   int    `json:"default_duration"`
				DefaultTrack      bool   `json:"default_track"`
				DisplayDimensions string `json:"display_dimensions"`
				EnabledTrack      bool   `json:"enabled_track"`
				ForcedTrack       bool   `json:"forced_track"`
				Language          string `json:"language"`
				TrackName         string `json:"track_name"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"tracks"`
	}

	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse mkvmerge output: %v", err)
	}

	info := &MediaInfo{
		FileName: result.FileName,
		Duration: float64(result.Container.Properties.Duration) / 1000000000, // наносекунды в секунды
	}

	for _, track := range result.Tracks {
		trackInfo := MediaInfoTrack{
			Number:            track.Properties.Number,
			DefaultDuration:   track.Properties.DefaultDuration,
			DefaultTrack:      track.Properties.DefaultTrack,
			DisplayDimensions: track.Properties.DisplayDimensions,
			EnabledTrack:      track.Properties.EnabledTrack,
			ForcedTrack:       track.Properties.ForcedTrack,
			Language:          track.Properties.Language,
			TrackName:         track.Properties.TrackName,
		}

		switch track.Type {
		case "video":
			info.VideoTracks = append(info.VideoTracks, trackInfo)
		case "audio":
			info.AudioTracks = append(info.AudioTracks, trackInfo)
		case "subtitles":
			info.Subtitles = append(info.Subtitles, trackInfo)
		}
	}

	return info, nil
}
