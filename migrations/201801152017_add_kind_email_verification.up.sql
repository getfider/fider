ALTER TABLE email_verifications ADD kind smallint NULL;

UPDATE email_verifications SET kind = 1 WHERE name IS NULL OR name = '';
UPDATE email_verifications SET kind = 2 WHERE name IS NOT NULL AND name != '';

ALTER TABLE email_verifications ALTER COLUMN kind SET NOT NULL;

ALTER TABLE email_verifications ADD user_id INT NULL REFERENCES users (id)