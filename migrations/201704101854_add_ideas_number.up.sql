ALTER TABLE ideas ADD NUMBER INT;

UPDATE ideas i
SET NUMBER = i2.seqnum
FROM (SELECT i2.*, row_number() OVER (PARTITION BY tenant_id ORDER BY created_on) AS seqnum FROM ideas i2) i2
WHERE i2.id = i.id;

ALTER TABLE ideas ALTER COLUMN NUMBER SET NOT NULL;