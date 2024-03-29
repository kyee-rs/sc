// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package postgresql

import (
	"context"
)

const createFile = `-- name: CreateFile :one
insert into files (id, name, mime, size, buffer, hash)
values ($1, $2, $3, $4, $5, $6)
returning id, name, mime, size, buffer, hash, created_at
`

type CreateFileParams struct {
	ID     string
	Name   string
	Mime   string
	Size   int64
	Buffer []byte
	Hash   string
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) (File, error) {
	row := q.db.QueryRowContext(ctx, createFile,
		arg.ID,
		arg.Name,
		arg.Mime,
		arg.Size,
		arg.Buffer,
		arg.Hash,
	)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Mime,
		&i.Size,
		&i.Buffer,
		&i.Hash,
		&i.CreatedAt,
	)
	return i, err
}

const getFile = `-- name: GetFile :one
select id, name, mime, size, buffer, hash, created_at
from files
where id = $1
limit 1
`

func (q *Queries) GetFile(ctx context.Context, id string) (File, error) {
	row := q.db.QueryRowContext(ctx, getFile, id)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Mime,
		&i.Size,
		&i.Buffer,
		&i.Hash,
		&i.CreatedAt,
	)
	return i, err
}

const getFileHash = `-- name: GetFileHash :one
select id, name, mime, size, buffer, hash, created_at
from files
where hash = $1
limit 1
`

func (q *Queries) GetFileHash(ctx context.Context, hash string) (File, error) {
	row := q.db.QueryRowContext(ctx, getFileHash, hash)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Mime,
		&i.Size,
		&i.Buffer,
		&i.Hash,
		&i.CreatedAt,
	)
	return i, err
}

const purgeFiles = `-- name: PurgeFiles :exec
delete
from files
where created_at < now() - interval '3 days'
`

func (q *Queries) PurgeFiles(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, purgeFiles)
	return err
}
