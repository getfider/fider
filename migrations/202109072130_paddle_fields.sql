alter table tenants_billing drop column stripe_customer_id;
alter table tenants_billing drop column stripe_subscription_id;
alter table tenants_billing drop column stripe_plan_id;

alter table tenants_billing add paddle_subscription_id varchar(255) not null;
alter table tenants_billing add paddle_plan_id varchar(255) not null;
alter table tenants_billing add status smallint not null;