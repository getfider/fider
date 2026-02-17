ALTER TABLE tenants ADD default_sort VARCHAR(50) NULL;
UPDATE tenants SET default_sort = 'most-wanted';
ALTER TABLE tenants ALTER COLUMN default_sort SET NOT NULL;
ALTER TABLE tenants ALTER COLUMN default_sort SET DEFAULT 'most-wanted';
