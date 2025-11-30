-- Remove Paddle and trial-related columns from tenants_billing
-- Keep paddle_subscription_id to identify customers who need to migrate to Stripe
-- Make paddle_subscription_id nullable (was NOT NULL from 202109072130 migration)

ALTER TABLE tenants_billing DROP COLUMN IF EXISTS paddle_plan_id;
ALTER TABLE tenants_billing DROP COLUMN IF EXISTS trial_ends_at;

-- Make paddle_subscription_id nullable so new Stripe subscriptions don't fail
ALTER TABLE tenants_billing ALTER COLUMN paddle_subscription_id DROP NOT NULL;

-- NOTE: paddle_subscription_id is intentionally kept for migration tracking
