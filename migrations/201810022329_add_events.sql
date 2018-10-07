CREATE TABLE IF NOT EXISTS events (
    id SERIAL NOT NULL,
    tenant_id INT NOT NULL,
    client_ip INET,
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);