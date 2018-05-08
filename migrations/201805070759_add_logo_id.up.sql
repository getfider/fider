CREATE UNIQUE INDEX upload_tenant_key ON uploads (id, tenant_id);

ALTER TABLE tenants ADD logo_id INT NULL;

ALTER TABLE tenants
   ADD CONSTRAINT tenants_logo_id_fkey
   FOREIGN KEY (logo_id, id) 
   REFERENCES uploads(id, tenant_id);
