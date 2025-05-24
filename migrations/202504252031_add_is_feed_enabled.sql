ALTER TABLE tenants ADD is_feed_enabled BOOLEAN NULL;

UPDATE tenants SET is_feed_enabled = true;

ALTER TABLE tenants ALTER COLUMN is_feed_enabled SET NOT NULL;