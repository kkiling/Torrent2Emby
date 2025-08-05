package sqlite

import (
	"context"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase/sqlitebase"
)

type Storage struct {
	base *sqlitebase.Storage
}

func (s *Storage) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	return s.base.RunTransaction(ctx, txFunc)
}

func NewStorage(config sqlitebase.Config, logger log.Logger) (*Storage, error) {
	s, err := sqlitebase.NewStorage(config, logger)
	if err != nil {
		return nil, err
	}
	return &Storage{
		base: s,
	}, nil
}

func NewTestStorage(base *sqlitebase.Storage) (*Storage, error) {
	return &Storage{
		base: base,
	}, nil
}
