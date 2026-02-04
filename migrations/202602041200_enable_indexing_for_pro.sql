-- Enable search engine indexing for all existing pro accounts
UPDATE tenants
SET prevent_indexing = false
WHERE is_pro = true;
