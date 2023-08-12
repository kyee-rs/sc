create table files
(
    id         text   not null primary key,
    name       text   not null,
    mime       text   not null,
    size       bigint not null,
    buffer     bytea  not null,
    hash       text   not null,
    created_at timestamptz default (now())
);

create unique index on files(hash);