package teststate

type CreateOptions struct {
	IdempotencyKey string
	Title          string
	Amount         int
}

func (c CreateOptions) GetIdempotencyKey() string {
	return c.IdempotencyKey
}

type Data struct {
	Counter int
	Title   string
	Amount  int
}

type WaitingInputOptions struct {
	IsComplete bool
	NewAmount  int
}

type StepType string

const (
	UnknownStep          StepType = "unknown"
	FirstStep            StepType = "first"
	TestErrorStep        StepType = "test_error"
	TestNoSaveChangeStep StepType = "no_save_change"
	WaitingInputStep     StepType = "waiting_input"
)

type Type string

const (
	TestType Type = "test_state"
)
