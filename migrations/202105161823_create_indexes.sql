CREATE INDEX IF NOT EXISTS post_tags_post_id_fkey ON post_tags (tenant_id,post_id);
CREATE INDEX IF NOT EXISTS comments_post_id_fkey ON comments (tenant_id,post_id);
CREATE INDEX IF NOT EXISTS post_votes_post_id_fkey ON post_votes (tenant_id,post_id);
CREATE INDEX IF NOT EXISTS post_subscribers_post_id_fkey ON post_subscribers (tenant_id,post_id);