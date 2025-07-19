ALTER TABLE tenants ADD prevent_indexing BOOLEAN NULL;

UPDATE tenants
SET
    prevent_indexing = false;

ALTER TABLE tenants
ALTER COLUMN prevent_indexing
SET
    NOT NULL;

ALTER TABLE tenants
ALTER COLUMN prevent_indexing
SET DEFAULT true;