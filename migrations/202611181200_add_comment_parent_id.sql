CREATE UNIQUE INDEX comments_id_tenant_key ON comments (id, tenant_id);

ALTER TABLE comments ADD COLUMN parent_id INT NULL;
ALTER TABLE comments ADD CONSTRAINT comments_parent_id_fkey
  FOREIGN KEY (parent_id, tenant_id) REFERENCES comments(id, tenant_id);

CREATE INDEX comments_parent_id_idx ON comments (tenant_id, parent_id) WHERE parent_id IS NOT NULL;
