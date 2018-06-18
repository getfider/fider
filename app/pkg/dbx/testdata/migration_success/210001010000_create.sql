create table dummy (
  id    int not null, 
  description  varchar(200) not null
);

insert into dummy (id, description) values (100, 'Description 100X');
insert into dummy (id, description) values (200, 'Description 200Y');
insert into dummy (id, description) values (300, 'Description 300Z');