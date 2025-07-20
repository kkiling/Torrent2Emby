package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kkiling/torrent2emby/internal/log"
	"github.com/kkiling/torrent2emby/internal/statemachine/storage"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type txKey struct{}

type Config struct {
	DSN string
}

type Storage struct {
	config Config
	db     *sql.DB
	logger log.Logger
}

func NewStorage(config Config, logger log.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Включаем поддержку внешних ключей
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed PRAGMA foreign_keys = ON;: %w", err)
	}

	return &Storage{config: config, db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNotFound
	}

	// Проверяем, является ли ошибка ошибкой уникальности
	sqliteErr := sqlite3.Error{}
	if ok := errors.As(err, &sqliteErr); ok {
		if strings.Contains(sqliteErr.Error(), "FOREIGN KEY constraint failed") {
			return err
		}
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return storage.ErrAlreadyExists
		}
	}

	return err
}

func (s *Storage) next(ctx context.Context) SQLExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		// Используем внутренний db из Tx
		return tx
	}
	return s.db
}

func (s *Storage) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	// Начало транзакции
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Отложенный вызов Rollback в случае ошибки
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				s.logger.Errorf("failed to rollback transaction: %v", err)
			}
		}
	}()

	ctxTx := context.WithValue(ctx, txKey{}, tx)
	if tErr := txFunc(ctxTx); tErr != nil {
		return fmt.Errorf("run transaction: %w", tErr)
	}

	// Фиксация транзакции
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
