create table if not exists comments (
     id          serial primary key,
     content     text null,
     idea_id     int not null,
     user_id     int not null,
     created_on  timestamptz not null,
     foreign key (idea_id) references ideas(id),
     foreign key (user_id) references users(id)
);