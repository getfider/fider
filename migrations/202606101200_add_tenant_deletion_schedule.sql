ALTER TABLE tenants ADD COLUMN scheduled_deletion_at TIMESTAMPTZ NULL;
ALTER TABLE tenants ADD COLUMN deletion_requested_by INT NULL;
ALTER TABLE tenants ADD COLUMN deletion_cancel_key   VARCHAR(64) NULL;
