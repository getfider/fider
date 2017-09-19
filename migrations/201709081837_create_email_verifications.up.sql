create table if not exists email_verifications (
     id             serial primary key,
     tenant_id      int not null,
     email          varchar(200) not null,
     created_on     timestamptz not null,
     key            varchar(32) not null,
     verified_on    timestamptz null,
     foreign key (tenant_id) references tenants(id)
);