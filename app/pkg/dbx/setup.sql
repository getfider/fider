TRUNCATE TABLE logs RESTART IDENTITY CASCADE;
TRUNCATE TABLE tenants RESTART IDENTITY CASCADE;

INSERT INTO tenants (name, subdomain, created_at, cname, invitation, welcome_message, status, is_private, custom_css) 
VALUES ('Demonstration', 'demo', now(), '', '', '', 1, false, '');

INSERT INTO users (name, email, tenant_id, created_at, role, status) 
VALUES ('Jon Snow', 'jon.snow@got.com', 1, now(), 3, 1);
INSERT INTO user_providers (user_id, tenant_id, provider, provider_uid, created_at) 
VALUES (1, 1, 'facebook', 'FB1234', now());

INSERT INTO users (name, email, tenant_id, created_at, role, status) 
VALUES ('Arya Stark', 'arya.stark@got.com', 1, now(), 1, 1);
INSERT INTO user_providers (user_id, tenant_id, provider, provider_uid, created_at) 
VALUES (2, 1, 'google', 'GO5678', now());

INSERT INTO users (name, email, tenant_id, created_at, role, status) 
VALUES ('Sansa Stark', 'sansa.stark@got.com', 1, now(), 1, 1);

INSERT INTO tenants (name, subdomain, created_at, cname, invitation, welcome_message, status, is_private, custom_css) 
VALUES ('Avengers', 'avengers', now(), 'feedback.avengers.com', '', '', 1, false, '');

INSERT INTO users (name, email, tenant_id, created_at, role, status) 
VALUES ('Tony Stark', 'tony.stark@avengers.com', 2, now(), 3, 1);
INSERT INTO user_providers (user_id, tenant_id, provider, provider_uid, created_at) 
VALUES (4, 2, 'facebook', 'FB2222', now());

INSERT INTO users (name, email, tenant_id, created_at, role, status) 
VALUES ('The Hulk', 'the.hulk@avengers.com', 2, now(), 1, 1);
INSERT INTO user_providers (user_id, tenant_id, provider, provider_uid, created_at) 
VALUES (5, 2, 'google', 'GO1111', now());
