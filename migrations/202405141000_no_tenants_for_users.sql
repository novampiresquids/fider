ALTER TABLE users DROP COLUMN tenant_id CASCADE;
ALTER TABLE users DROP COLUMN role CASCADE;
ALTER TABLE user_settings DROP COLUMN tenant_id CASCADE;

CREATE TABLE members (
    user_id     int NOT NULL,
    tenant_id   int NOT NULL,
    role        int NOT NULL,

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)