-- Migrate from billing status to simple plan system
-- Active status (2) becomes Pro plan, all others become Free plan

-- Add is_pro column to tenants table
ALTER TABLE tenants ADD COLUMN is_pro BOOLEAN NOT NULL DEFAULT false;

-- Migrate existing data: Active billing status (2) -> Pro plan
UPDATE tenants t
SET is_pro = true
FROM tenants_billing tb
WHERE t.id = tb.tenant_id
AND tb.status = 2;

-- Drop unused columns from tenants_billing
ALTER TABLE tenants_billing DROP COLUMN IF EXISTS status;
ALTER TABLE tenants_billing DROP COLUMN IF EXISTS subscription_ends_at;
