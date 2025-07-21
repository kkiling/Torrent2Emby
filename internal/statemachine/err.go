package statemachine

import "errors"

var (
	// ErrNotFound не найден
	ErrNotFound = errors.New("state not found")
	// ErrAlreadyExists уже существует
	ErrAlreadyExists = errors.New("state already exists")
	// ErrInTerminalStatus уже в терминальном статусе
	ErrInTerminalStatus = errors.New("state already in terminal status")
)
