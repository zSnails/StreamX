// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const allMedia = `-- name: AllMedia :many
SELECT id, hash, title, creator FROM media
`

func (q *Queries) AllMedia(ctx context.Context) ([]Medium, error) {
	rows, err := q.db.Query(ctx, allMedia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Medium
	for rows.Next() {
		var i Medium
		if err := rows.Scan(
			&i.ID,
			&i.Hash,
			&i.Title,
			&i.Creator,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createMedia = `-- name: CreateMedia :exec
INSERT INTO media (hash, title, creator) VALUES ($1, $2, $3)
`

type CreateMediaParams struct {
	Hash    string
	Title   string
	Creator string
}

func (q *Queries) CreateMedia(ctx context.Context, arg CreateMediaParams) error {
	_, err := q.db.Exec(ctx, createMedia, arg.Hash, arg.Title, arg.Creator)
	return err
}

const findMedia = `-- name: FindMedia :many
SELECT id, hash, title, creator FROM media WHERE SIMILARITY(title, $1) > 0.1
`

func (q *Queries) FindMedia(ctx context.Context, similarity string) ([]Medium, error) {
	rows, err := q.db.Query(ctx, findMedia, similarity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Medium
	for rows.Next() {
		var i Medium
		if err := rows.Scan(
			&i.ID,
			&i.Hash,
			&i.Title,
			&i.Creator,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStoredMedia = `-- name: GetStoredMedia :one
SELECT name, fileoid FROM media_files WHERE name = $1
`

func (q *Queries) GetStoredMedia(ctx context.Context, name pgtype.Text) (MediaFile, error) {
	row := q.db.QueryRow(ctx, getStoredMedia, name)
	var i MediaFile
	err := row.Scan(&i.Name, &i.Fileoid)
	return i, err
}

const mediaExists = `-- name: MediaExists :one
SELECT exists (SELECT 1 FROM media WHERE hash = $1 LIMIT 1)
`

func (q *Queries) MediaExists(ctx context.Context, hash string) (bool, error) {
	row := q.db.QueryRow(ctx, mediaExists, hash)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const storeMedia = `-- name: StoreMedia :exec
INSERT INTO media_files (name, fileoid) values ($1, $2)
`

type StoreMediaParams struct {
	Name    pgtype.Text
	Fileoid pgtype.Uint32
}

func (q *Queries) StoreMedia(ctx context.Context, arg StoreMediaParams) error {
	_, err := q.db.Exec(ctx, storeMedia, arg.Name, arg.Fileoid)
	return err
}
