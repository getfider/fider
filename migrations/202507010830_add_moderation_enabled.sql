ALTER TABLE tenants ADD is_moderation_enabled BOOLEAN NULL;

UPDATE tenants SET is_moderation_enabled = false;

ALTER TABLE tenants ALTER COLUMN is_moderation_enabled SET NOT NULL;