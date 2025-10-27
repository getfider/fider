CREATE TABLE tenant_providers (
    id           SERIAL PRIMARY KEY,
    tenant_id    INT NOT NULL,
    provider     VARCHAR(40) NOT NULL,
    is_enabled   BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, provider),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX idx_tenant_providers_tenant_id ON tenant_providers(tenant_id);
