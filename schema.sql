CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    hash VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    creator VARCHAR(255) NOT NULL
);

CREATE INDEX idx_media_title_gin ON media USING GIN (title gin_trgm_ops);
