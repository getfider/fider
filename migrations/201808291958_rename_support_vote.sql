ALTER TABLE post_supporters RENAME TO post_votes;
ALTER INDEX post_supporters_pkey RENAME TO post_votes_pkey;
ALTER TABLE post_votes RENAME CONSTRAINT post_supporters_post_id_fkey TO post_votes_post_id_fkey;
ALTER TABLE post_votes RENAME CONSTRAINT post_supporters_tenant_id_fkey TO post_votes_tenant_id_fkey;
ALTER TABLE post_votes RENAME CONSTRAINT post_supporters_user_id_fkey TO post_votes_user_id_fkey;

UPDATE tenants 
SET custom_css = replace(custom_css, '.c-support-counter', '.c-vote-counter')
WHERE custom_css LIKE '%.c-support-counter%';

UPDATE tenants 
SET custom_css = replace(custom_css, '.m-supported', '.m-voted')
WHERE custom_css LIKE '%.m-supported%';