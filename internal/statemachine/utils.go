package statemachine

import (
	"time"

	"github.com/google/uuid"
)

type uuidGenerator struct{}

func (uuidGenerator) New() uuid.UUID {
	return uuid.New()
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}
