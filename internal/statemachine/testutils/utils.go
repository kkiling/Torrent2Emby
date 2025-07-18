package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// findUp ищет файл с именем `filename` в текущей и родительских директориях.
// Возвращает полный путь к файлу или пустую строку, если файл не найден.
func findUp(filename string) string {
	dir, err := os.Getwd() // Текущая директория
	if err != nil {
		return ""
	}

	for {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path // Файл найден
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Достигли корневой директории
		}
		dir = parent // Переходим в родительскую директорию
	}

	return "" // Файл не найден
}

var envPath = ""

func init() {
	envPath = findUp(".testenv")
	if envPath == "" {
		panic(".testenv not found in current or parent directories")
	}

	// Загружаем переменные из найденного файла
	if err := godotenv.Load(envPath); err != nil {
		panic(err)
	}
}

func GetSqliteTestDNS(t *testing.T) string {
	// Загружаем переменные окружения из .testenv файла
	err := godotenv.Overload(envPath)
	require.NoError(t, err, "Failed to load .testenv file")

	return os.Getenv("SQLITE_DSN")
}
