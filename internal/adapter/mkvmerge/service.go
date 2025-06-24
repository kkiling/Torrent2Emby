package mkvmerge

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Merge(params MergeParams) error {
	// Проверка существования основного видеофайла
	if _, err := os.Stat(params.VideoInputFile); os.IsNotExist(err) {
		return fmt.Errorf("input video file does not exist: %s", params.VideoInputFile)
	}

	// Проверка аудиодорожек
	for _, track := range params.AudioTracks {
		if _, err := os.Stat(track.Path); os.IsNotExist(err) {
			return fmt.Errorf("audio track file does not exist: %s", track.Path)
		}
	}

	// Проверка субтитров
	for _, track := range params.SubtitleTracks {
		if _, err := os.Stat(track.Path); os.IsNotExist(err) {
			return fmt.Errorf("subtitle file does not exist: %s", track.Path)
		}
	}

	args := []string{"-o", params.VideoOutputFile, params.VideoInputFile}

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
	fmt.Println("Executing command:", "mkvmerge", strings.Join(args, " "))

	// Создаем команду
	cmd := exec.Command("mkvmerge", args...)

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
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Читаем вывод в реальном времени
	go scanOutput(stdoutPipe)
	go scanOutput(stderrPipe)

	// Ждем завершения
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("mkvmerge failed: %v", err)
	}

	return nil
}

func scanOutput(reader io.Reader) {
	buf := make([]byte, 1024)
	var leftover []byte

	for {
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
					fmt.Printf("%s\n", string(line))
				}
			}
		}

		if err != nil {
			if err != io.EOF {
				fmt.Printf("read error: %v\n", err)
			}
			break
		}
	}

	// Выводим оставшиеся данные
	if len(leftover) > 0 {
		fmt.Printf("%s\n", string(leftover))
	}
}
