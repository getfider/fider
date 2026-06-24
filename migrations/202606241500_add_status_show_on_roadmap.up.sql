-- Per-status opt-in flag for the Roadmap page. The Roadmap previously
-- inferred lanes from kind=active OR kind=closed-completed; this column lets
-- admins explicitly publish any status (including kind=open or custom kinds)
-- to the roadmap.
--
-- Backfill: seeded rows whose kind=active or closed-completed default to true
-- so the existing roadmap stays unchanged for tenants that haven't customised.
-- Everything else (open, closed-declined, duplicate, custom slugs) defaults to
-- false until an admin flips it on.

ALTER TABLE statuses ADD COLUMN IF NOT EXISTS show_on_roadmap BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE statuses
SET show_on_roadmap = TRUE
WHERE kind IN ('active', 'closed-completed');
