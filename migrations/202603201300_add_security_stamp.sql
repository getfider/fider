ALTER TABLE users ADD COLUMN security_stamp VARCHAR(64) NOT NULL DEFAULT '';
UPDATE users SET security_stamp = md5(random()::text || id::text || tenant_id::text);
