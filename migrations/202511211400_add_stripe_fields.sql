ALTER TABLE tenants_billing ADD COLUMN stripe_customer_id VARCHAR(255) NULL;
ALTER TABLE tenants_billing ADD COLUMN stripe_subscription_id VARCHAR(255) NULL;
