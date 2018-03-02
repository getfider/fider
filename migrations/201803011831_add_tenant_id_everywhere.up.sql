CREATE UNIQUE INDEX idea_id_tenant_id_key ON ideas (tenant_id, id);
CREATE UNIQUE INDEX user_id_tenant_id_key ON users (tenant_id, id);
CREATE UNIQUE INDEX tag_id_tenant_id_key ON tags (tenant_id, id);



---- TABLE: comments

-- add tenant_id
ALTER TABLE comments ADD tenant_id INT NULL;
UPDATE comments SET tenant_id = (SELECT tenant_id FROM ideas WHERE ideas.id = comments.idea_id);
ALTER TABLE comments ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE comments
   ADD CONSTRAINT comments_tenant_id_fkey
   FOREIGN KEY (tenant_id) 
   REFERENCES tenants(id);
   
-- comments <-> ideas
ALTER TABLE comments
    DROP CONSTRAINT comments_idea_id_fkey RESTRICT;

ALTER TABLE comments
   ADD CONSTRAINT comments_idea_id_fkey
   FOREIGN KEY (idea_id, tenant_id) 
   REFERENCES ideas(id, tenant_id);
   
-- comments <-> users
ALTER TABLE comments
    DROP CONSTRAINT comments_user_id_fkey RESTRICT;

ALTER TABLE comments
   ADD CONSTRAINT comments_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);



---- TABLE: email_verifications

-- email_verifications <-> users
ALTER TABLE email_verifications
    DROP CONSTRAINT email_verifications_user_id_fkey RESTRICT;
    
ALTER TABLE email_verifications
   ADD CONSTRAINT email_verifications_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);



---- TABLE: user_providers

-- add tenant_id
ALTER TABLE user_providers ADD tenant_id INT NULL;
UPDATE user_providers SET tenant_id = (SELECT tenant_id FROM users WHERE users.id = user_providers.user_id);
ALTER TABLE user_providers ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE user_providers
   ADD CONSTRAINT user_providers_tenant_id_fkey
   FOREIGN KEY (tenant_id) 
   REFERENCES tenants(id);

-- user_providers <-> users
ALTER TABLE user_providers
    DROP CONSTRAINT user_providers_user_id_fkey RESTRICT;
    
ALTER TABLE user_providers
   ADD CONSTRAINT user_providers_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);



---- TABLE: user_settings

-- add tenant_id
ALTER TABLE user_settings ADD tenant_id INT NULL;
UPDATE user_settings SET tenant_id = (SELECT tenant_id FROM users WHERE users.id = user_settings.user_id);
ALTER TABLE user_settings ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE user_settings
   ADD CONSTRAINT user_settings_tenant_id_fkey
   FOREIGN KEY (tenant_id) 
   REFERENCES tenants(id);

-- user_settings <-> users
ALTER TABLE user_settings
    DROP CONSTRAINT user_settings_user_id_fkey RESTRICT;
    
ALTER TABLE user_settings
   ADD CONSTRAINT user_settings_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);



---- TABLE: idea_subscribers

-- add tenant_id
ALTER TABLE idea_subscribers ADD tenant_id INT NULL;
UPDATE idea_subscribers SET tenant_id = (SELECT tenant_id FROM users WHERE users.id = idea_subscribers.user_id);
ALTER TABLE idea_subscribers ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE idea_subscribers
   ADD CONSTRAINT idea_subscribers_tenant_id_fkey
   FOREIGN KEY (tenant_id) 
   REFERENCES tenants(id);
   
-- idea_subscribers <-> ideas
ALTER TABLE idea_subscribers
    DROP CONSTRAINT idea_subscribers_idea_id_fkey RESTRICT;

ALTER TABLE idea_subscribers
   ADD CONSTRAINT idea_subscribers_idea_id_fkey
   FOREIGN KEY (idea_id, tenant_id) 
   REFERENCES ideas(id, tenant_id);
   
-- idea_subscribers <-> users
ALTER TABLE idea_subscribers
    DROP CONSTRAINT idea_subscribers_user_id_fkey RESTRICT;

ALTER TABLE idea_subscribers
   ADD CONSTRAINT idea_subscribers_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);



---- TABLE: idea_supporters

-- add tenant_id
ALTER TABLE idea_supporters ADD tenant_id INT NULL;
UPDATE idea_supporters SET tenant_id = (SELECT tenant_id FROM users WHERE users.id = idea_supporters.user_id);
ALTER TABLE idea_supporters ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE idea_supporters
   ADD CONSTRAINT idea_supporters_tenant_id_fkey
   FOREIGN KEY (tenant_id) 
   REFERENCES tenants(id);
   
-- idea_supporters <-> ideas
ALTER TABLE idea_supporters
    DROP CONSTRAINT idea_supporters_idea_id_fkey RESTRICT;

ALTER TABLE idea_supporters
   ADD CONSTRAINT idea_supporters_idea_id_fkey
   FOREIGN KEY (idea_id, tenant_id) 
   REFERENCES ideas(id, tenant_id);
   
-- idea_supporters <-> users
ALTER TABLE idea_supporters
    DROP CONSTRAINT idea_supporters_user_id_fkey RESTRICT;

ALTER TABLE idea_supporters
   ADD CONSTRAINT idea_supporters_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);

---- TABLE: idea_tags

-- add tenant_id
ALTER TABLE idea_tags ADD tenant_id INT NULL;
UPDATE idea_tags SET tenant_id = (SELECT tenant_id FROM tags WHERE tags.id = idea_tags.tag_id);
ALTER TABLE idea_tags ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE idea_tags
   ADD CONSTRAINT idea_tags_tenant_id_fkey
   FOREIGN KEY (tenant_id) 
   REFERENCES tenants(id);
   
-- idea_tags <-> ideas
ALTER TABLE idea_tags
    DROP CONSTRAINT idea_tags_idea_id_fkey RESTRICT;

ALTER TABLE idea_tags
   ADD CONSTRAINT idea_tags_idea_id_fkey
   FOREIGN KEY (idea_id, tenant_id) 
   REFERENCES ideas(id, tenant_id);
   
-- idea_tags <-> users
ALTER TABLE idea_tags
    DROP CONSTRAINT idea_tags_created_by_id_fkey RESTRICT;

ALTER TABLE idea_tags
   ADD CONSTRAINT idea_tags_created_by_id_fkey
   FOREIGN KEY (created_by_id, tenant_id) 
   REFERENCES users(id, tenant_id);
   
-- idea_tags <-> tags
ALTER TABLE idea_tags
    DROP CONSTRAINT idea_tags_tag_id_fkey RESTRICT;

ALTER TABLE idea_tags
   ADD CONSTRAINT idea_tags_tag_id_fkey
   FOREIGN KEY (tag_id, tenant_id) 
   REFERENCES tags(id, tenant_id);

---- TABLE: notifications
   
-- notifications <-> ideas
ALTER TABLE notifications
    DROP CONSTRAINT notifications_idea_id_fkey RESTRICT;

ALTER TABLE notifications
   ADD CONSTRAINT notifications_idea_id_fkey
   FOREIGN KEY (idea_id, tenant_id) 
   REFERENCES ideas(id, tenant_id);
   
-- notifications <-> users
ALTER TABLE notifications
    DROP CONSTRAINT notifications_user_id_fkey RESTRICT;

ALTER TABLE notifications
   ADD CONSTRAINT notifications_user_id_fkey
   FOREIGN KEY (user_id, tenant_id) 
   REFERENCES users(id, tenant_id);
   
-- notifications <-> users (author)
ALTER TABLE notifications
    DROP CONSTRAINT notifications_author_id_fkey RESTRICT;

ALTER TABLE notifications
   ADD CONSTRAINT notifications_author_id_fkey
   FOREIGN KEY (author_id, tenant_id) 
   REFERENCES users(id, tenant_id);