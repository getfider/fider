-- Add post_id to mention_notifications table
ALTER TABLE mention_notifications
ADD COLUMN post_id INTEGER;

-- Add foreign key constraint on post_id that allows null values
ALTER TABLE mention_notifications ADD CONSTRAINT mention_notifications_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts (id);

ALTER TABLE mention_notifications
DROP CONSTRAINT IF EXISTS mention_notifications_comment_id_fkey;

ALTER TABLE mention_notifications
DROP CONSTRAINT IF EXISTS unique_mention_notification;

-- Add a new foreign key constraint that allows null values
ALTER TABLE mention_notifications ADD CONSTRAINT mention_notifications_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES comments (id);

-- Update the unique constraint to include post_id
ALTER TABLE mention_notifications
DROP CONSTRAINT IF EXISTS mention_notifications_tenant_id_user_id_comment_id_key;

-- Add a check constraint to ensure at least one of comment_id or post_id is not null
ALTER TABLE mention_notifications ADD CONSTRAINT mention_notifications_comment_or_post_check CHECK (
    (comment_id IS NOT NULL)
    OR (post_id IS NOT NULL)
);