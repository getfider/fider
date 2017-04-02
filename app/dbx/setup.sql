TRUNCATE TABLE tenants RESTART IDENTITY CASCADE;

INSERT INTO tenants (id, name, subdomain, created_on) VALUES (300, 'Demonstration', 'demo', now());

INSERT INTO users (id, name, email, tenant_id, created_on) VALUES (300, 'Jon Snow', 'jon.snow@got.com', 300, now());
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (300, 'facebook', 'FB1234', now());

INSERT INTO users (id, name, email, tenant_id, created_on) VALUES (301, 'Arya Stark', 'arya.stark@got.com', 300, now());
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (301, 'google', 'GO5678', now());

INSERT INTO tenants (id, name, subdomain, created_on, cname) VALUES (400, 'Orange Inc.', 'orange', now(), 'feedback.orangeinc.com');

INSERT INTO users (id, name, email, tenant_id, created_on) VALUES (400, 'Tony Stark', 'tony.stark@avengers.com', 400, now());
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (400, 'facebook', 'FB2222', now());

INSERT INTO users (id, name, email, tenant_id, created_on) VALUES (401, 'The Hulk', 'the.hulk@avengers.com', 400, now());
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (401, 'google', 'GO1111', now());
