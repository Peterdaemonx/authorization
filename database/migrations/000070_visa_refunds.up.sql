ALTER TABLE refunds ADD COLUMN authorization_id_response STRING(6);
ALTER TABLE refunds ADD COLUMN retrieval_reference_number STRING(12);

CREATE TABLE visa_refunds
(
    refund_id              STRING(36) NOT NULL,
    created_at             TIMESTAMP NOT NULL,
    transaction_identifier INT64
) PRIMARY KEY(refund_id),
INTERLEAVE IN PARENT refunds ON DELETE CASCADE;