create table blocks
(
    id              bigserial                 not null
        constraint blocks_pk
            primary key,
    created         timestamptz default now() not null,
    updated         timestamptz default now() not null,
    height          numeric                   not null check ( height > 0 ) unique,
    hash            text                      not null check ( hash != '' ) unique,
    block_time  timestamptz not null
);

insert into blocks(height,hash,block_time) values (657957, 'B02EF50091031EF9AAC7E8BBDD98395B89BC42EB90B9138EFE2E5DF79186EC7B','2023-01-12T19:22:57.245474096Z');

---- create above / drop below ----
drop table blocks;
