alter table contracts add column closed_height bigint not null default 0;
---- create above / drop below ----
drop view open_contracts_v;
alter table contracts drop column closed_height;
