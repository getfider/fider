-- Add license key storage to tenants_billing table
ALTER TABLE tenants_billing ADD COLUMN license_key TEXT;
