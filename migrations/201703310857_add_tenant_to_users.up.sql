ALTER TABLE users ADD tenant_id INT REFERENCES tenants(id);

UPDATE users
SET tenant_id = ideas.tenant_id
FROM ideas
WHERE ideas.user_id = users.id;

UPDATE users
SET tenant_id = (SELECT id FROM tenants LIMIT 1)
WHERE tenant_id IS NULL