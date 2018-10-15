ALTER TABLE comments ADD deleted_at TIMESTAMPTZ NULL;
ALTER TABLE comments ADD deleted_by_id INT NULL;

ALTER TABLE comments
   ADD CONSTRAINT comments_deleted_by_id_fkey
   FOREIGN KEY (deleted_by_id, tenant_id) 
   REFERENCES users(id, tenant_id);