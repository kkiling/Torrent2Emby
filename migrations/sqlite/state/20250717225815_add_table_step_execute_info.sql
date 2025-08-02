-- +goose Up
-- +goose StatementBegin
CREATE TABLE step_execute_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state_id BLOB NOT NULL,
    start_executed_at TIMESTAMP NOT NULL,
    complete_executed_at TIMESTAMP,
    error TEXT,
    preview_step TEXT NOT NULL,
    next_step TEXT,
    FOREIGN KEY (state_id) REFERENCES state(id) ON DELETE CASCADE
);

CREATE INDEX idx_step_execute_state_id ON step_execute_info(state_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE step_execute_info;
-- +goose StatementEnd
