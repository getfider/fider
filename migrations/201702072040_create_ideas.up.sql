create table if not exists ideas (
     id           serial primary key,
     title        varchar(100) not null,
     description  text null,
     tenant_id    int not null,
     created_on   timestamptz not null default now(),
     modified_on  timestamptz not null default now(),
     foreign key (tenant_id) references tenants(id)
);