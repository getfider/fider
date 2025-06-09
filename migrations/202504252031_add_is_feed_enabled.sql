ALTER TABLE tenants ADD is_feed_enabled BOOLEAN NULL;

UPDATE tenants
SET is_feed_enabled = CASE
    WHEN is_private = true THEN false
    ELSE true
END;

ALTER TABLE tenants ALTER COLUMN is_feed_enabled SET NOT NULL;