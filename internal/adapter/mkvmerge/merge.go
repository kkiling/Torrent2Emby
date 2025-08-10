package mkvmerge

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kkiling/goplatform/log"
)

type Merge struct {
	logger log.Logger
}

func NewMerge(logger log.Logger) *Merge {
	return &Merge{
		logger: logger.Named("mkvmerge"),
	}
}

func (s *Merge) Merge(ctx context.Context, params MergeParams, outputChan chan<- OutputMessage) error {
	info, err := s.GetMediaInfo(params.VideoInputFile)
	if err != nil {
		return fmt.Errorf("get media info: %w", err)
	}

	// Проверка существования основного видеофайла
	if _, err = os.Stat(params.VideoInputFile); os.IsNotExist(err) {
		return fmt.Errorf("input video file does not exist: %s", params.VideoInputFile)
	}

	// Проверка аудиодорожек
	for _, track := range params.AudioTracks {
		if _, err = os.Stat(track.Path); os.IsNotExist(err) {
			return fmt.Errorf("audio track file does not exist: %s", track.Path)
		}
	}

	// Проверка субтитров
	for _, track := range params.SubtitleTracks {
		if _, err = os.Stat(track.Path); os.IsNotExist(err) {
			return fmt.Errorf("subtitle file does not exist: %s", track.Path)
		}
	}

	args := []string{"-o", params.VideoOutputFile}

	// Снимаем default со всех старых аудио в исходном файле
	for _, audio := range info.AudioTracks {
		// audio.Number — это номер трека в контейнере (1-based), mkvmerge ждёт 0-based
		trackID := audio.Number - 1
		args = append(args, "--default-track", fmt.Sprintf("%d:no", trackID))
	}

	args = append(args, params.VideoInputFile)

	// Добавляем аудиодорожки
	for _, track := range params.AudioTracks {
		if track.Language != "" {
			args = append(args, "--language", "0:"+track.Language)
		}
		args = append(args,
			"--track-name", "0:"+track.Name,
			"--default-track", fmt.Sprintf("0:%v", track.Default),
			track.Path, // Путь к файлу идет ПОСЛЕ флагов!
		)
	}

	// Добавляем субтитры
	for _, track := range params.SubtitleTracks {
		if track.Language != "" {
			args = append(args, "--language", "0:"+track.Language)
		}
		args = append(args,
			"--track-name", "0:"+track.Name,
			"--default-track", fmt.Sprintf("0:%v", track.Default),
			track.Path, // Путь к файлу идет ПОСЛЕ флагов!
		)
	}

	// Для отладки
	debugMsg := "mkvmerge " + strings.Join(args, " ")
	outputChan <- OutputMessage{Type: InfoMessageType, Content: debugMsg}

	// Создаем команду
	cmd := exec.CommandContext(ctx, "mkvmerge", args...)
	// Настраиваем пайпы
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Запускаем команду
	if errStart := cmd.Start(); errStart != nil {
		return fmt.Errorf("error starting command: %v", errStart)
	}

	// Читаем вывод в реальном времени и отправляем в канал
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		s.scanOutput(ctx, stdoutPipe, outputChan, InfoMessageType)
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		s.scanOutput(ctx, stderrPipe, outputChan, ErrorMessageType)
	}()

	// Ждем завершения
	if errWait := cmd.Wait(); errWait != nil {
		return fmt.Errorf("mkvmerge failed: %w", errWait)
	}
	wg.Wait()

	return nil
}

// Измененная функция scanOutput теперь является методом Service и принимает канал
func (s *Merge) scanOutput(ctx context.Context, reader io.Reader, outputChan chan<- OutputMessage, msgType MessageType) {
	buf := make([]byte, 1024)
	var leftover []byte

	for ctx.Err() == nil {
		n, err := reader.Read(buf)
		if n > 0 {
			data := append(leftover, buf[:n]...)
			lines := bytes.Split(data, []byte{'\r'})

			// Последний элемент может быть неполной строкой
			for i, line := range lines {
				if i == len(lines)-1 {
					leftover = line
					continue
				}

				// Обрабатываем только непустые строки
				if len(line) > 0 {
					outputChan <- OutputMessage{
						Type:    msgType,
						Content: string(line),
					}
				}
			}
		}

		if err != nil {
			if err != io.EOF {
				outputChan <- OutputMessage{
					Type:    ErrorMessageType,
					Content: fmt.Sprintf("read error: %v", err),
				}
			}
			break
		}
	}

	// Выводим оставшиеся данные
	if len(leftover) > 0 {
		outputChan <- OutputMessage{
			Type:    msgType,
			Content: string(leftover),
		}
	}
}
