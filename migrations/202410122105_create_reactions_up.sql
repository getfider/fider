create table if not exists reactions (
     id          serial primary key,
     emoji       varchar(8) not null,
     comment_id  int not null,
     user_id     int not null,
     created_on  timestamptz not null,
     foreign key (comment_id) references comments(id),
     foreign key (user_id) references users(id)
);

ALTER TABLE reactions ADD CONSTRAINT unique_reaction UNIQUE (comment_id, user_id, emoji);
