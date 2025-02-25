ALTER TABLE tenants ADD COLUMN profanity_words TEXT NULL;

UPDATE tenants SET profanity_words = '';

ALTER TABLE tenants ALTER COLUMN profanity_words SET NOT NULL;