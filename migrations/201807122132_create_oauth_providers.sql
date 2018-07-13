create table if not exists oauth_providers (
  id                   serial not null, 
  tenant_id            int not null,
  logo_id              int null,
  provider             varchar(30) not null,
  display_name         varchar(50) not null,
  status               int not null,
  client_id            varchar(100) not null,
  client_secret        varchar(500) not null,
  authorize_url        varchar(300) not null,
  token_url            varchar(300) not null,
  profile_url          varchar(300) not null,
  scope                varchar(100) not null,
  json_user_id_path    varchar(100) not null,
  json_user_name_path  varchar(100) not null,
  json_user_email_path varchar(100) not null,
  created_on           timestamptz not null default now(),
  primary key (id),
  foreign key (tenant_id) references tenants(id),
  foreign key (logo_id, tenant_id) references uploads(id, tenant_id)
);

CREATE UNIQUE INDEX tenant_id_provider_key ON oauth_providers (tenant_id, provider);