CREATE TABLE api_consumers_permissions (
    permission_id STRING(36) NOT NULL,
    psp_id STRING(36) NOT NULL,
    api_consumer_id STRING(36) NOT NULL,
    CONSTRAINT FK_permissions FOREIGN KEY (permission_id) REFERENCES permissions (permission_id),
) PRIMARY KEY(psp_id, api_consumer_id, permission_id),
INTERLEAVE IN PARENT api_consumers ON DELETE CASCADE;