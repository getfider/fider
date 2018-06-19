create table if not exists logs (
  id            serial not null, 
  tag           varchar(50) not null,
  level         varchar(50) not null,
  text          text not null,
  properties    jsonb null,
  created_on    timestamptz not null default now(),
  primary key (id)
);