-- +goose Up
-- +goose StatementBegin
-- Таблица для хранения результатов слияния
CREATE TABLE IF NOT EXISTS mkv_merge (
    id UUID PRIMARY KEY,
    idempotency_key TEXT NOT NULL,  -- uuid.UUID
    params JSONB NOT NULL,
    status INTEGER NOT NULL,
    error TEXT,
    created_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    progress REAL
);

CREATE UNIQUE INDEX idx_mkv_merge_idempotency_key ON mkv_merge(idempotency_key);

-- Таблица для логов процесса слияния
CREATE TABLE IF NOT EXISTS mkv_merge_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    merge_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    type INTEGER NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (merge_id) REFERENCES mkv_merge(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE mkv_merge;
DROP TABLE mkv_merge_logs;
-- +goose StatementEnd
