create table if not exists users (
     id           serial primary key,
     name         varchar(100) null,
     email        varchar(200) not null,
     created_on   timestamptz not null default now(),
     modified_on  timestamptz not null default now()
);

create table if not exists user_providers (
     user_id      int not null,
     provider     varchar(40) not null,
     provider_uid varchar(100) not null,
     created_on   timestamptz not null default now(),
     modified_on  timestamptz not null default now(),
     primary key (user_id, provider),
     foreign key (user_id) references users(id)
);