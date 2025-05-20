package config

import (
	"errors"
	"os"
)

// getEnv возвращает значение переменной окружения с проверкой на пустое значение
func getEnv[T any](key string, convert func(string) (T, error)) (T, error) {
	var zero T
	value := os.Getenv(key)
	if value == "" {
		return zero, errors.New(key + " is not set")
	}

	result, err := convert(value)
	if err != nil {
		return zero, errors.New("failed to parse " + key + ": " + err.Error())
	}

	return result, nil
}

// getEnvString упрощенная версия для строк
func getEnvString(key string) (string, error) {
	return getEnv(key, func(s string) (string, error) { return s, nil })
}
