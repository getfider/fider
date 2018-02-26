create table if not exists notifications (
  id          serial not null, 
  tenant_id   int not null,
  user_id     int not null,
  title       varchar(160) not null,
  link        varchar(2048) null,
  read        boolean not null, 
  idea_id     int not null,
  author_id   int not null,
  created_on  timestamptz not null default now(),
  updated_on  timestamptz not null default now(),
  primary key (id),
  foreign key (tenant_id) references tenants(id),
  foreign key (user_id) references users(id),
  foreign key (author_id) references users(id),
  foreign key (idea_id) references ideas(id)
);