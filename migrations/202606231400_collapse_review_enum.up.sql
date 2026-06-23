-- Collapse the HCM-only PostReview=7 legacy enum value. Any post that was
-- previously responded as "review" gets its status_slug filled in (defensive,
-- the dual-write should have already done this) and its int status reset to 0
-- (PostOpen). Identity continues to live in posts.status_slug; the upcoming
-- migration drops the int column entirely.
--
-- Safe on databases that never had enum 7 rows — both UPDATEs are no-ops.

UPDATE posts SET status_slug = 'review' WHERE status = 7 AND (status_slug IS NULL OR status_slug = '');
UPDATE posts SET status = 0 WHERE status = 7;
