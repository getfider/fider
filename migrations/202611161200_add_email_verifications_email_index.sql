CREATE INDEX IF NOT EXISTS email_verifications_tenant_email_idx ON email_verifications (tenant_id, email, kind, created_at);
