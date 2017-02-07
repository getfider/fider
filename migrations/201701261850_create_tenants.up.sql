create table if not exists tenants (
     id           serial primary key,
     name         varchar(60) not null,
     domain       varchar(40) not null,
     created_on   timestamptz not null default now(),
     modified_on  timestamptz not null default now()
);