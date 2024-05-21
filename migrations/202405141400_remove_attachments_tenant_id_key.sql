ALTER TABLE attachments DROP CONSTRAINT attachments_tenant_id_fkey;
ALTER TABLE blobs DROP CONSTRAINT blobs_tenant_id_fkey;
ALTER TABLE email_verifications DROP CONSTRAINT email_verifications_tenant_id_fkey;
DROP INDEX tenant_subdomain_unique_idx;
CREATE UNIQUE INDEX tenant_user_unique_idx ON members (tenant_id, user_id);
ALTER TABLE user_providers DROP COLUMN tenant_id CASCADE;