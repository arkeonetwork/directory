create table provider_metadata
(
    id                            bigserial                 not null
        constraint provider_metadata_pk
            primary key,
    created                       timestamptz default now() not null,
    updated                       timestamptz default now() not null,
    provider_id                   bigint                    not null references providers (id),
    version                       text                      not null,
    moniker                       text,
    website                       text,
    description                   text,
    location                      text,
    port                          text,
    proxy_host                    text,
    source_chain                  text,
    event_stream_host             text,
    claim_store_location          text,
    free_rate_limit               bigint,
    free_rate_limit_duration      bigint,
    subscribe_rate_limit          bigint,
    subscribe_rate_limit_duration bigint,
    paygo_rate_limit              bigint,
    paygo_rate_limit_duration     bigint
);

alter table provider_metadata
    add constraint prov_version_uniq unique (provider_id, version);

---- create above / drop below ----
drop table provider_metadata;
