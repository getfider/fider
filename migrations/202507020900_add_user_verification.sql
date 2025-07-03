-- Add user verification column to users table
-- This allows users to be verified for automatic approval of their content
ALTER TABLE users ADD COLUMN is_verified boolean NOT NULL DEFAULT false;

-- Add index for performance
CREATE INDEX idx_users_verification ON users(tenant_id, is_verified) WHERE is_verified = true;