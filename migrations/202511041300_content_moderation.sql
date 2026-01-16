-- Add is_moderation_enabled column to tenants table
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS is_moderation_enabled BOOLEAN NULL;

UPDATE tenants SET is_moderation_enabled = false WHERE is_moderation_enabled IS NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'tenants'
        AND column_name = 'is_moderation_enabled'
        AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE tenants ALTER COLUMN is_moderation_enabled SET NOT NULL;
    END IF;
END $$;

-- Add is_approved column to posts table
ALTER TABLE posts ADD COLUMN IF NOT EXISTS is_approved BOOLEAN NULL;

UPDATE posts SET is_approved = true WHERE is_approved IS NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'posts'
        AND column_name = 'is_approved'
        AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE posts ALTER COLUMN is_approved SET NOT NULL;
    END IF;
END $$;

-- Add is_approved column to comments table
ALTER TABLE comments ADD COLUMN IF NOT EXISTS is_approved BOOLEAN NULL;

UPDATE comments SET is_approved = true WHERE is_approved IS NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'comments'
        AND column_name = 'is_approved'
        AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE comments ALTER COLUMN is_approved SET NOT NULL;
    END IF;
END $$;

-- Add user trust column to users table
-- This allows users to be trusted for automatic approval of their content
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_trusted boolean NOT NULL DEFAULT false;

-- Add index for performance
CREATE INDEX IF NOT EXISTS idx_users_trust ON users(tenant_id, is_trusted) WHERE is_trusted = true;