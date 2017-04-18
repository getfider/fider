ALTER TABLE ideas ADD supporters INT;

UPDATE ideas SET supporters = 0;

ALTER TABLE ideas ALTER COLUMN supporters SET NOT NULL;

CREATE TABLE IF NOT EXISTS idea_supporters (
     user_id     int not null,
     idea_id     int not null,
     created_on  timestamptz not null,
     primary key (user_id, idea_id),
     foreign key (idea_id) references ideas(id),
     foreign key (user_id) references users(id)
);