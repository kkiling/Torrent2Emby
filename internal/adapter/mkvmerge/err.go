package mkvmerge

import "errors"

var (
	// ErrNotFound не найден
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists уже существует
	ErrAlreadyExists = errors.New("already exists")
)
