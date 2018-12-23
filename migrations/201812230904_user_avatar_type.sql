alter table users add column avatar_type smallint null;
alter table users add column avatar_bkey varchar(512) null;

update users set avatar_type = 2; -- gravatar

alter table users alter column avatar_type set not null;