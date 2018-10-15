ALTER TABLE comments ADD status INT NULL;

UPDATE comments SET status = 1;

ALTER TABLE comments ALTER COLUMN status SET NOT NULL;