-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS video_content (
    -- Идентификатор записи (uuid.UUID как строка)
    id TEXT PRIMARY KEY NOT NULL,
    -- Дата создания 
    created_at TIMESTAMP NOT NULL,
    -- ContentID: вариант "фильм"
    movie_id INTEGER,
    -- ContentID: вариант "сериал"
    tvshow_id INTEGER,
    -- Номер сезона сериала
    season_number INTEGER,
    -- Статус доставки
    delivery_status TEXT NOT NULL,
    -- Массив состояний State в JSON:
    -- [{ "state_id": "<uuid>", "type": "<value>" }, ...]
    states_json TEXT,
    -- Ограничения
    CONSTRAINT delivery_status_chk CHECK (
        delivery_status IN ('failed', 'in_progress', 'delivered' , 'updating' , 'deleting' , 'deleted')
    ),
    -- Ровно один вариант ContentID: либо movie, либо tvshow+season
    CONSTRAINT content_choice_chk CHECK (
        (movie_id IS NOT NULL AND tvshow_id IS NULL AND season_number IS NULL)
            OR
        (movie_id IS NULL AND tvshow_id IS NOT NULL AND season_number IS NOT NULL)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS video_content;
-- +goose StatementEnd
