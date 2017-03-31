alter table user_providers drop column modified_on;
alter table user_providers alter column created_on drop default;