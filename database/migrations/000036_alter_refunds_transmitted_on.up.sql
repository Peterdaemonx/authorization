ALTER TABLE refunds DROP COLUMN transmitted_on;
ALTER TABLE refunds ADD COLUMN transmitted_at TIMESTAMP;
ALTER TABLE refunds ALTER COLUMN transmitted_at TIMESTAMP NOT NULL;