CREATE INDEX IF NOT EXISTS comments_edited_at
ON comments (tenant_id, COALESCE(edited_at, created_at) ASC)
WHERE deleted_at IS NULL;
