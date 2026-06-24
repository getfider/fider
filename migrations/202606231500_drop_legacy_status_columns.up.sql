-- Phase 2 cleanup: posts.status_slug is now the sole identifier for a post's
-- status. Drop the legacy int column on posts, the transitional join column
-- on statuses, and tighten posts.status_slug to NOT NULL.
--
-- Defensive backfill first in case any pre-202606231300 row slipped through
-- without a slug. Maps int 0..5 to the built-in slugs the schema seeds; rows
-- with status not in that range (impossible after 202606231400) fall back to
-- 'open'.

UPDATE posts SET status_slug = CASE status
    WHEN 0 THEN 'open'
    WHEN 1 THEN 'started'
    WHEN 2 THEN 'completed'
    WHEN 3 THEN 'declined'
    WHEN 4 THEN 'planned'
    WHEN 5 THEN 'duplicate'
    WHEN 6 THEN 'deleted'
    ELSE 'open'
END
WHERE status_slug IS NULL OR status_slug = '';

ALTER TABLE posts ALTER COLUMN status_slug SET NOT NULL;
ALTER TABLE posts DROP COLUMN status;

ALTER TABLE statuses DROP COLUMN legacy_enum;
