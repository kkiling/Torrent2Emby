package teststate

import "github.com/google/uuid"

type CreateOptions struct {
	IdempotencyKey uuid.UUID
	Title          string
	Priority       int
}

func (c CreateOptions) GetIdempotencyKey() uuid.UUID {
	return c.IdempotencyKey
}

type Data struct {
	ID    int
	Title string
}

type Status string

const (
	PendingStatus Status = "pending"
	DoneStatus    Status = "done"
)
