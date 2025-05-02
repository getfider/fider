CREATE TABLE
    if not exists draft_posts (
        id serial primary key,
        title varchar(100) not null,
        code varchar(12) not null,
        description text null,
        created_at timestamptz not null default now ()
    );

CREATE TABLE
    if not exists draft_attachments (
        id serial not null,
        draft_post_id int not null,
        attachment_bkey varchar(512) not null,
        primary key (id),
        foreign key (draft_post_id) references draft_posts (id)
    );