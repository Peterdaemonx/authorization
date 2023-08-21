CREATE TABLE users_permissions
(
    permission_id STRING(36) NOT NULL,
    person_guid   STRING(36) NOT NULL,
    CONSTRAINT FK_user_permissions FOREIGN KEY (permission_id) REFERENCES permissions (permission_id),
) PRIMARY KEY(person_guid, permission_id),
INTERLEAVE IN PARENT users ON DELETE CASCADE;