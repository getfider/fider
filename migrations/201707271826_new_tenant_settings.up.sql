alter table tenants add invitation varchar(100) null;
alter table tenants add welcome_message text null;

update tenants set invitation = '', welcome_message = '';