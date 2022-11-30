create table provider_status
(
    id     serial not null
        constraint provider_status_pk
            primary key,
    status text   not null unique
);

insert into provider_status(status) values ('Online');
insert into provider_status(status) values ('Offline');

create table contract_types
(
    id  serial not null
        constraint contract_types_pk
            primary key,
    val text   not null unique
);

insert into contract_types(val)
values ('PayAsYouGo');
insert into contract_types(val)
values ('Subscription');

---- create above / drop below ----
-- undo --
drop table contract_types;
drop table provider_status;
