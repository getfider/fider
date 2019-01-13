create table if not exists attachments (
  id              serial not null,
  tenant_id       int not null,
  post_id         int not null,
  comment_id      int null,
  user_id         int not null,
  attachment_bkey varchar(512) not null,
  primary key (id),
  foreign key (tenant_id) references tenants(id),
  foreign key (post_id) references posts(id),
  foreign key (user_id) references users(id),
  foreign key (comment_id) references comments(id)
);