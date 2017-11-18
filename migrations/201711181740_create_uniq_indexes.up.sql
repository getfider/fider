CREATE UNIQUE INDEX idea_number_tenant_key ON ideas (tenant_id, number);
CREATE UNIQUE INDEX tag_slug_tenant_key ON tags (tenant_id, slug);