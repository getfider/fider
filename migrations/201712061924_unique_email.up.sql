CREATE UNIQUE INDEX user_email_unique_idx ON users (tenant_id, email) WHERE email != '';
CREATE UNIQUE INDEX user_provider_unique_idx ON user_providers (user_id, provider);
CREATE UNIQUE INDEX tenant_subdomain_unique_idx ON tenants (subdomain);
CREATE UNIQUE INDEX tenant_cname_unique_idx ON tenants (cname) WHERE cname != '';