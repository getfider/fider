alter table tenants drop column modified_on;
alter table tenants alter column created_on drop default;