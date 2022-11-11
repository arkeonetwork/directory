-- drop table if exists providers;

create table providers
(
    id      bigserial                 not null
        constraint providers_pk
            primary key,
    created timestamptz default now() not null,
    updated timestamptz default now() not null,
    pubkey  text                      not null,
    chain   text                      not null,
    bond    numeric                   not null
);

alter table providers
    add constraint pubkey_chain_uniq unique (pubkey, chain);


