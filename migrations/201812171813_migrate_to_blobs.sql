insert into blobs (key, tenant_id, size, content_type, file, created_at, modified_at)
select 
    'logos/'||md5(cast(id as varchar)) || '.' ||  split_part(content_type, '/', 2) as key,
    tenant_id,
    size,
    content_type,
    file,
    created_at,
    created_at as modified_at
from uploads;

alter table tenants add column logo_bkey varchar(512) null;
alter table oauth_providers add column logo_bkey varchar(512) null;

update tenants
set logo_bkey = 'logos/'||md5(cast(u.id as varchar)) || '.' ||  split_part(u.content_type, '/', 2)
from  uploads u
where u.id = tenants.logo_id;

update oauth_providers
set logo_bkey = 'logos/'||md5(cast(u.id as varchar)) || '.' ||  split_part(u.content_type, '/', 2)
from  uploads u
where u.id = oauth_providers.logo_id;

alter table tenants drop column logo_id;
alter table oauth_providers drop column logo_id;

update tenants set logo_bkey = '' where logo_bkey is null;
update oauth_providers set logo_bkey = '' where logo_bkey is null;

alter table tenants alter column logo_bkey SET NOT NULL;
alter table oauth_providers alter column logo_bkey SET NOT NULL;

drop table uploads;