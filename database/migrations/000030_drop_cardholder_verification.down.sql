ALTER TABLE authorizations ADD COLUMN cardholder_verification STRING(30);
ALTER TABLE authorizations ALTER COLUMN cardholder_verification STRING(30) NOT NULL;

ALTER TABLE refunds ADD COLUMN cardholder_verification STRING(30);
ALTER TABLE refunds ALTER COLUMN cardholder_verification STRING(30) NOT NULL;