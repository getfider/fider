create table if not exists uploads (
  id            serial not null, 
  tenant_id     int not null,
  size          int not null,
  content_type  varchar(200) not null,
  file          bytea not null,
  created_on    timestamptz not null default now(),
  primary key (id),
  foreign key (tenant_id) references tenants(id)
);