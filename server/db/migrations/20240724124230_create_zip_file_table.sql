-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS zip_file (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    s3_key TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE zip_file IF EXISTS;
-- +goose StatementEnd
