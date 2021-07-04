ALTER TABLE tenants ADD locale VARCHAR(10) NULL;
UPDATE tenants SET locale = 'en';
ALTER TABLE tenants ALTER COLUMN locale SET NOT NULL;