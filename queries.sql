-- name: FindMedia :many
SELECT * FROM media WHERE title ILIKE $1;

-- name: AllMedia :many
SELECT * FROM media;

-- name: CreateMedia :exec
INSERT INTO media (hash, title, creator) VALUES ($1, $2, $3);

--SELECT * FROM media WHERE SIMILARITY(title, $1) > 0.3;
