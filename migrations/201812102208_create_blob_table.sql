create table if not exists blobs (
  id            serial not null,
  key           varchar(512) not null, 
  tenant_id     int null,
  size          bigint not null,
  content_type  varchar(200) not null,
  file          bytea not null,
  created_at    timestamptz not null default now(),
  modified_at   timestamptz not null default now(),
  primary key (id),
  unique (tenant_id, key),
  foreign key (tenant_id) references tenants(id)
);