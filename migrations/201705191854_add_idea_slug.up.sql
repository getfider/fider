ALTER TABLE ideas ADD slug varchar(100);

UPDATE ideas SET slug = '';

ALTER TABLE ideas ALTER COLUMN slug SET NOT NULL;