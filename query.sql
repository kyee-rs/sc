-- name: GetFile :one
select *
from files
where id = $1
limit 1;

-- name: GetFileHash :one
select *
from files
where hash = $1
limit 1;

-- name: CreateFile :one
insert into files (id, name, mime, size, buffer, hash)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: PurgeFiles :exec
select drop_chunks('files', older_than := interval '3d');