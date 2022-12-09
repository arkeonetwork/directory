alter table contracts add column closed_height bigint not null default 0;
---- create above / drop below ----
alter table contracts drop column closed_height;
