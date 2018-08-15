UPDATE user_settings
SET key = 'event_notification_new_post'
WHERE key = 'event_notification_new_idea';

ALTER TABLE idea_subscribers RENAME COLUMN idea_id TO post_id;
ALTER TABLE idea_supporters RENAME COLUMN idea_id TO post_id;
ALTER TABLE idea_tags RENAME COLUMN idea_id TO post_id;
ALTER TABLE comments RENAME COLUMN idea_id TO post_id;
ALTER TABLE notifications RENAME COLUMN idea_id TO post_id;

ALTER TABLE ideas RENAME TO posts;
ALTER TABLE idea_subscribers RENAME TO post_subscribers;
ALTER TABLE idea_supporters RENAME TO post_supporters;
ALTER TABLE idea_tags RENAME TO post_tags;

ALTER INDEX idea_id_tenant_id_key RENAME TO post_id_tenant_id_key;
ALTER INDEX idea_number_tenant_key RENAME TO post_number_tenant_key;
ALTER INDEX idea_slug_tenant_key RENAME TO post_slug_tenant_key;
ALTER INDEX idea_subscribers_pkey RENAME TO post_subscribers_pkey;
ALTER INDEX idea_supporters_pkey RENAME TO post_supporters_pkey;
ALTER INDEX idea_tags_pkey RENAME TO post_tags_pkey;
ALTER INDEX ideas_pkey RENAME TO posts_pkey;

ALTER TABLE comments RENAME CONSTRAINT comments_idea_id_fkey TO comments_post_id_fkey;
ALTER TABLE post_subscribers RENAME CONSTRAINT idea_subscribers_idea_id_fkey TO post_subscribers_post_id_fkey;
ALTER TABLE post_subscribers RENAME CONSTRAINT idea_subscribers_tenant_id_fkey TO post_subscribers_tenant_id_fkey;
ALTER TABLE post_subscribers RENAME CONSTRAINT idea_subscribers_user_id_fkey TO post_subscribers_user_id_fkey;
ALTER TABLE post_supporters RENAME CONSTRAINT idea_supporters_idea_id_fkey TO post_supporters_post_id_fkey;
ALTER TABLE post_supporters RENAME CONSTRAINT idea_supporters_tenant_id_fkey TO post_supporters_tenant_id_fkey;
ALTER TABLE post_supporters RENAME CONSTRAINT idea_supporters_user_id_fkey TO post_supporters_user_id_fkey;
ALTER TABLE post_tags RENAME CONSTRAINT idea_tags_created_by_id_fkey TO post_tags_created_by_id_fkey;
ALTER TABLE post_tags RENAME CONSTRAINT idea_tags_idea_id_fkey TO post_tags_post_id_fkey;
ALTER TABLE post_tags RENAME CONSTRAINT idea_tags_tag_id_fkey TO post_tags_tag_id_fkey;
ALTER TABLE post_tags RENAME CONSTRAINT idea_tags_tenant_id_fkey TO post_tags_tenant_id_fkey;
ALTER TABLE posts RENAME CONSTRAINT ideas_original_id_fkey TO posts_original_id_fkey;
ALTER TABLE posts RENAME CONSTRAINT ideas_tenant_id_fkey TO posts_tenant_id_fkey;
ALTER TABLE posts RENAME CONSTRAINT ideas_user_id_fkey TO posts_user_id_fkey;
ALTER TABLE notifications RENAME CONSTRAINT notifications_idea_id_fkey TO notifications_post_id_fkey;

ALTER SEQUENCE ideas_id_seq RENAME TO posts_id_seq;