-- Remove Paddle and trial-related columns from tenants_billing
-- We're moving to Stripe-only billing with no trial periods

ALTER TABLE tenants_billing DROP COLUMN IF EXISTS paddle_subscription_id;
ALTER TABLE tenants_billing DROP COLUMN IF EXISTS paddle_plan_id;
ALTER TABLE tenants_billing DROP COLUMN IF EXISTS trial_ends_at;
