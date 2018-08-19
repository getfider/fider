ALTER TABLE users ADD api_key VARCHAR(32) NULL;
ALTER TABLE users ADD api_key_date TIMESTAMPTZ NULL;

CREATE UNIQUE INDEX users_api_key ON users (tenant_id, api_key);