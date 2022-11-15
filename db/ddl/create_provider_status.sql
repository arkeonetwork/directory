drop table if exists provider_status;

create table provider_status
(
    id     serial not null
        constraint provider_status_pk
            primary key,
    status text   not null unique
);

insert into provider_status(status) values ('Online');
insert into provider_status(status) values ('Offline');
