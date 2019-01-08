DROP INDEX post_slug_tenant_key;
CREATE UNIQUE INDEX post_slug_tenant_key ON posts (tenant_id, slug) WHERE status <> 6;