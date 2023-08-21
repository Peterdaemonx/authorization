ALTER TABLE refunds DROP COLUMN mc_response_code;
ALTER TABLE refunds ADD COLUMN mc_response_code STRING(20);