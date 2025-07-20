-- +goose Up
-- +goose StatementBegin

-- Таблица для хранения изображений
CREATE TABLE images (
    id TEXT PRIMARY KEY,
    w342 TEXT,
    original TEXT
);

-- Таблица для хранения сезонов
CREATE TABLE seasons (
    id INTEGER PRIMARY KEY,
    tv_show_id INTEGER NOT NULL,
    air_date TIMESTAMP NOT NULL,  -- time.Time
    episode_count INTEGER,
    name TEXT,
    overview TEXT,
    poster_id TEXT,
    season_number INTEGER,
    vote_average REAL,
    FOREIGN KEY (poster_id) REFERENCES images(id),
    FOREIGN KEY (tv_show_id) REFERENCES tv_shows(id)
);

-- Таблица для хранения TVShow
CREATE TABLE tv_shows (
    id INTEGER PRIMARY KEY,
    name TEXT,
    original_name TEXT,
    overview TEXT,
    poster_id TEXT,
    first_air_date TIMESTAMP NOT NULL,  -- time.Time
    vote_average REAL,
    vote_count INTEGER,
    popularity REAL,
    backdrop_id INTEGER,
    genres TEXT, -- массив genres в формате JSON
    last_air_date TIMESTAMP NOT NULL,  -- time.Time
    next_episode_to_air TIMESTAMP NOT NULL,  -- time.Time
    number_of_episodes INTEGER,
    number_of_seasons INTEGER,
    origin_country TEXT, -- массив origin_country в формате JSON
    status TEXT,
    tagline TEXT,
    type TEXT,
    FOREIGN KEY (poster_id) REFERENCES images(id),
    FOREIGN KEY (backdrop_id) REFERENCES images(id)
);

CREATE TABLE episodes (
    id INTEGER PRIMARY KEY,
    season_id INTEGER NOT NULL,
    air_date TIMESTAMP NOT NULL,
    episode_number INTEGER NOT NULL,
    episode_type TEXT,
    name TEXT,
    overview TEXT,
    runtime INTEGER,
    still_id TEXT,
    vote_average REAL,
    vote_count INTEGER,
    FOREIGN KEY (season_id) REFERENCES seasons(id),
    FOREIGN KEY (still_id) REFERENCES images(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tv_shows;
DROP TABLE seasons;
DROP TABLE images;
-- +goose StatementEnd
