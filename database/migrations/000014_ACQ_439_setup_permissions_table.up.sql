CREATE TABLE permissions
(
    permission_id STRING(36) NOT NULL,
    code          STRING(100) NOT NULL,
    label         STRING(100) NOT NULL,
) PRIMARY KEY (permission_id);