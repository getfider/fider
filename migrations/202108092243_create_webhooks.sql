CREATE TABLE IF NOT EXISTS webhooks (
  id SERIAL PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  type SMALLINT NOT NULL,
  status SMALLINT NOT NULL,
  url TEXT NOT NULL,
  content TEXT NULL,
  http_method VARCHAR(50) NOT NULL,
  http_headers JSONB NULL,
  tenant_id INT NOT NULL,
  FOREIGN KEY (tenant_id) REFERENCES tenants (id)
);