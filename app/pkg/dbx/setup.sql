TRUNCATE TABLE tenants RESTART IDENTITY CASCADE;

INSERT INTO tenants (name, subdomain, created_on, cname, invitation, welcome_message, status) VALUES ('Demonstration', 'demo', now(), '', '', '', 1);

INSERT INTO users (name, email, tenant_id, created_on, role) VALUES ('Jon Snow', 'jon.snow@got.com', 1, now(), 3);
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (1, 'facebook', 'FB1234', now());

INSERT INTO users (name, email, tenant_id, created_on, role) VALUES ('Arya Stark', 'arya.stark@got.com', 1, now(), 1);
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (2, 'google', 'GO5678', now());

INSERT INTO tenants (name, subdomain, created_on, cname, invitation, welcome_message, status) VALUES ('Avengers', 'avengers', now(), 'feedback.avengers.com', '', '', 1);

INSERT INTO users (name, email, tenant_id, created_on, role) VALUES ('Tony Stark', 'tony.stark@avengers.com', 2, now(), 3);
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (3, 'facebook', 'FB2222', now());

INSERT INTO users (name, email, tenant_id, created_on, role) VALUES ('The Hulk', 'the.hulk@avengers.com', 2, now(), 1);
INSERT INTO user_providers (user_id, provider, provider_uid, created_on) VALUES (4, 'google', 'GO1111', now());
