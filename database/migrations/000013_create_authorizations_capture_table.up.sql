CREATE TABLE authorization_captures
(
    authorization_id STRING(36) NOT NULL,
    capture_id STRING(36) NOT NULL,
    amount INT64 NOT NULL,
    currency STRING(3) NOT NULL,
    is_final BOOL NOT NULL,
    status STRING(10) NOT NULL,
    created_at TIMESTAMP NOT NULL,
) PRIMARY KEY(authorization_id, capture_id),
INTERLEAVE IN PARENT authorizations ON DELETE CASCADE;