CREATE TABLE IF NOT EXISTS demo (
     id   integer PRIMARY KEY DEFAULT nextval('demo'),
     name varchar(40) NOT NULL CHECK (name <> '')
);