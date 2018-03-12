-- add tenant_id
ALTER TABLE comments ADD edited_on TIMESTAMPTZ NULL;
ALTER TABLE comments ADD edited_by_id INT NULL;

ALTER TABLE comments
   ADD CONSTRAINT comments_edited_by_id_fkey
   FOREIGN KEY (edited_by_id, tenant_id) 
   REFERENCES users(id, tenant_id);