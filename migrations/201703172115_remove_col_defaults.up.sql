alter table ideas drop column modified_on;
alter table ideas alter column created_on drop default;

alter table users drop column modified_on;
alter table users alter column created_on drop default;