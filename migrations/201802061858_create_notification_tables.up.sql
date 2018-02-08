create table if not exists idea_subscribers (
  user_id     int not null,
  idea_id     int not null,
  created_on  timestamptz not null default now(),
  updated_on  timestamptz not null default now(),
  status      smallint not null,
  primary key (user_id, idea_id),
  foreign key (idea_id) references ideas(id),
  foreign key (user_id) references users(id)
);

create table if not exists user_settings (
  id          serial primary key,
  user_id     int not null,
  key         varchar(100) not null,
  value       varchar(100) null,
  foreign key (user_id) references users(id)
);

create unique index user_settings_uq_key on user_settings (user_id, key);