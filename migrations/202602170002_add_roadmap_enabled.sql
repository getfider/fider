ALTER TABLE tenants ADD is_roadmap_enabled BOOLEAN NULL;
UPDATE tenants SET is_roadmap_enabled = true;
ALTER TABLE tenants ALTER COLUMN is_roadmap_enabled SET NOT NULL;
ALTER TABLE tenants ALTER COLUMN is_roadmap_enabled SET DEFAULT true;
