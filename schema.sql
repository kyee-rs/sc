create table files
(
    created_at timestamptz default (now()),
    id         text   not null,
    name       text   not null,
    mime       text   not null,
    size       bigint not null,
    buffer     bytea  not null,
    hash       text   not null
);

select create_hypertable('files', 'created_at', migrate_data := true, chunk_time_interval := interval '3d');