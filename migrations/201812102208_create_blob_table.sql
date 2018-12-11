create table if not exists blobs (
  key           varchar(512) not null, 
  tenant_id     int null,
  size          int not null,
  content_type  varchar(200) not null,
  file          bytea not null,
  created_on    timestamptz not null default now(),
  primary key (id),
  foreign key (tenant_id) references tenants(id)
);