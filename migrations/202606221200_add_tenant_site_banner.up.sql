-- Site-wide banner controlled from Site Settings. Off by default for
-- existing installs, so this migration adds no visible behavior change.
ALTER TABLE tenants ADD site_banner_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE tenants ADD site_banner_message TEXT NOT NULL DEFAULT '';
ALTER TABLE tenants ADD site_banner_variant VARCHAR(20) NOT NULL DEFAULT 'info';
