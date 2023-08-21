ALTER TABLE refunds DROP COLUMN transmitted_at;
ALTER TABLE refunds ADD COLUMN transmitted_on DATE;