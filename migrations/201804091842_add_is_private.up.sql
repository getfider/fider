ALTER TABLE tenants ADD is_private BOOLEAN NULL;

UPDATE tenants SET is_private = false;

ALTER TABLE tenants ALTER COLUMN is_private SET NOT NULL;