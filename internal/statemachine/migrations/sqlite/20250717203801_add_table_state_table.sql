-- +goose Up
-- +goose StatementBegin
CREATE TABLE state (
    id BLOB PRIMARY KEY,  -- uuid.UUID
    idempotency_key TEXT NOT NULL,  -- uuid.UUID
    created_at TIMESTAMP NOT NULL,  -- time.Time
    updated_at TIMESTAMP NOT NULL,  -- time.Time
    status INTEGER NOT NULL,  -- uint8
    step TEXT NOT NULL,  -- string
    type TEXT NOT NULL,  -- string
    data BLOB,  -- []byte (JSONB equivalent)
    fail_data BLOB,  -- []byte (JSONB equivalent)
    meta_data BLOB  -- []byte (JSONB equivalent)
);

CREATE UNIQUE INDEX idx_state_idempotency_key ON state(idempotency_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE state;
-- +goose StatementEnd
