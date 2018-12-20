CREATE UNIQUE INDEX blobs_unique_global_key ON blobs (key, tenant_id) WHERE tenant_id IS NOT NULL;
CREATE UNIQUE INDEX blobs_unique_tenant_key ON blobs (key) WHERE tenant_id IS NULL;

insert into blobs (key, tenant_id, size, content_type, file, created_at, modified_at)
select 'autocert/'||key, null, octet_length(data), 'application/x-pem-file', data, created_at, created_at from autocert_cache