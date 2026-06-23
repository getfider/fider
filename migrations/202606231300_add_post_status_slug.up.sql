-- Phase 2 of feedback.fider.io/posts/111: make the post status text the source
-- of truth so admin-defined custom statuses can be assigned to posts.
--
-- posts.status (int) stays as fallback for any code path that still consults
-- the legacy enum. New writes go through status_slug; reads prefer it.

ALTER TABLE posts ADD status_slug VARCHAR(50) NULL;

-- Backfill from the per-tenant statuses table via legacy_enum mapping.
-- Built-in statuses (Open / Started / Completed / Declined / Planned /
-- Duplicate) all have legacy_enum set; any tenant that also seeded Review
-- (legacy_enum=7) gets resolved here too.
UPDATE posts p
SET status_slug = s.slug
FROM statuses s
WHERE s.tenant_id = p.tenant_id
  AND s.legacy_enum = p.status;

CREATE INDEX posts_tenant_status_slug_idx ON posts (tenant_id, status_slug);
