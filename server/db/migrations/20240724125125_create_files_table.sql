-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS extracted_files (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    s3_key VARCHAR(255) NOT NULL,
    zip_id INTEGER NOT NULL,
    content TEXT,
    searchable_content TSVECTOR,
    file_size BIGINT NOT NULL, 
    mime_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_extracted_files_searchable_content ON extracted_files USING GIN (searchable_content);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files IF EXISTS;
DROP INDEX idx_files_searchable_content IF EXISTS;
-- +goose StatementEnd
