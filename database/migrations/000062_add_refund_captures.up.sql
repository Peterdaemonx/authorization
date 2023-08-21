CREATE TABLE refund_captures
(
    refund_id STRING(36) NOT NULL,
    capture_id STRING(36) NOT NULL,
    reference STRING(100),
    amount INT64 NOT NULL,
    currency STRING(3) NOT NULL,
    is_final BOOL NOT NULL,
    status STRING(10) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
) PRIMARY KEY(refund_id, capture_id),
INTERLEAVE IN PARENT refunds ON DELETE CASCADE;