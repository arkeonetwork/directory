drop table if exists contract_types;

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
