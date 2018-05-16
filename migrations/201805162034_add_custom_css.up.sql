ALTER TABLE tenants ADD custom_css TEXT NULL;

UPDATE tenants SET custom_css = '';

ALTER TABLE tenants ALTER COLUMN custom_css SET NOT NULL;