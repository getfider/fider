ALTER TABLE oauth_providers ADD is_trusted BOOLEAN default false;

CREATE UNIQUE INDEX oauth_provider_uq ON oauth_providers (provider);