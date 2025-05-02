CREATE TABLE
    if not exists draft_posts (
        id serial primary key,
        title varchar(100) not null,
        uuid not null default gen_random_uuid (),
        description text null,
        created_on timestamptz not null default now (),
    );