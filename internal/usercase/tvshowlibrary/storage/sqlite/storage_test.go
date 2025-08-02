package sqlite

import (
	"testing"

	"github.com/kkiling/goplatform/log"
	"github.com/stretchr/testify/require"

	"github.com/kkiling/torrent2emby/internal/statemachine/testutils"
)

func setupTestDB(t *testing.T) *Storage {
	// Инициализируем хранилище
	cfg := Config{
		DSN: testutils.GetSqliteTestDNS(t), // Берем DSN из переменных окружения
	}
	logger := log.NewLogger(log.DebugLevel)

	s, err := NewStorage(cfg, logger)
	require.NoError(t, err)

	// Возвращаем хранилище и функцию очистки
	return s
}
