-- Add welcome_header column to tenants table
ALTER TABLE tenants ADD COLUMN welcome_header TEXT NULL DEFAULT '';
