ALTER TABLE tenants ADD COLUMN block_disposable_emails BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE email_domain_rules (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL,
    rule_type VARCHAR(10) NOT NULL CHECK (rule_type IN ('deny', 'allow')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, domain, rule_type)
);

CREATE INDEX email_domain_rules_tenant_idx ON email_domain_rules (tenant_id, rule_type);
