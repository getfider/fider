ALTER TABLE ideas ADD status INT;

UPDATE ideas SET status = 0;

ALTER TABLE ideas ALTER COLUMN status SET NOT NULL;