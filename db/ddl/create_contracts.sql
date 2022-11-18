drop table if exists contracts;

create table contracts
(
    id              bigserial                 not null
        constraint contracts_pk
            primary key,
    created         timestamptz default now() not null,
    updated         timestamptz default now() not null,
    provider_id     bigint                    not null references providers (id),
    delegate_pubkey text                      not null check ( delegate_pubkey != '' ),
    client_pubkey   text                      not null check ( client_pubkey != '' ),
    height          bigint                    not null check ( height > 0 ),
    contract_type   text                      not null references contract_types (val),
    duration        bigint                    not null,
    rate            bigint                    not null,
    open_cost       bigint                    not null
);

alter table contracts
    add constraint pubkey_prov_dlgt_uniq unique (provider_id, delegate_pubkey);

-- may be good'nuf with unique constraint
-- create index contracts_prov_id_idx on provider_mod_events (provider_id);
