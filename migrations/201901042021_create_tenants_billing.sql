create table if not exists tenants_billing (
  tenant_id              int not null,
  trial_ends_at          timestamptz not null,
  subscription_ends_at   timestamptz null,
  stripe_customer_id     varchar(255) null,
  stripe_subscription_id varchar(255) null,
  stripe_plan_id         varchar(255) null,
  primary key (tenant_id),
  foreign key (tenant_id) references tenants(id)
);