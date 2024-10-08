-- name: FindMedia :many
SELECT * FROM media WHERE SIMILARITY(title, $1) > 0.1;

-- name: AllMedia :many
SELECT * FROM media;

-- name: CreateMedia :exec
INSERT INTO media (hash, title, creator) VALUES ($1, $2, $3);

-- name: MediaExists :one
SELECT exists (SELECT 1 FROM media WHERE hash = $1 LIMIT 1);

-- name: StoreMedia :exec
INSERT INTO media_files (name, fileoid) values ($1, $2);

-- name: GetStoredMedia :one
SELECT * FROM media_files WHERE name = $1;
