create table
     if not exists mention_notifications (
          id serial not null,
          tenant_id int not null,
          user_id int not null,
          comment_id int null,
          created_on timestamptz not null default now (),
          primary key (id),
          foreign key (tenant_id) references tenants (id),
          foreign key (user_id) references users (id),
          foreign key (comment_id) references comments (id),
          constraint unique_mention_notification unique (tenant_id, user_id, comment_id)
     );

create index idx_mention_notifications_tenant_user on notification_logs (tenant_id, user_id);