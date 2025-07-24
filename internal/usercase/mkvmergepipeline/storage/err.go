package storage

import "errors"

var (
	// ErrNotFound объект не найден в базе
	ErrNotFound = errors.New("entity not found")
	// ErrAlreadyExists запись уже существует
	ErrAlreadyExists = errors.New("entity already exists")
)
