package storage

import "context"

// Transactor интерфейс транзакций
type Transactor interface {
	RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error
}
