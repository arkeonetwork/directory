drop table if exists indexer_status;

create table indexer_status
(
    id          serial not null
        constraint indexer_status_pk
            primary key,
    created     timestamptz default now() not null,
    updated     timestamptz default now() not null,
    height      numeric not null check ( height >= 0 )
);