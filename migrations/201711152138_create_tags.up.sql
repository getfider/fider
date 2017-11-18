create table if not exists tags (
     id             serial primary key,
     tenant_id      int not null,
     name           varchar(30) not null,
     slug           varchar(30) not null,
     color          varchar(6) not null,
     is_public      boolean not null,
     created_on     timestamptz not null,
     foreign key (tenant_id) references tenants(id)
);

create table if not exists idea_tags (
     tag_id         int not null,
     idea_id        int not null,
     created_on     timestamptz not null,
     created_by_id  int not null,
     primary key (tag_id, idea_id),
     foreign key (idea_id) references ideas(id),
     foreign key (tag_id) references tags(id),
     foreign key (created_by_id) references users(id)
);