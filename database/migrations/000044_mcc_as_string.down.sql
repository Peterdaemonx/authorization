ALTER TABLE authorizations DROP COLUMN merchant_category_code;
ALTER TABLE authorizations ADD COLUMN merchant_category_code INT64;

ALTER TABLE refunds DROP COLUMN merchant_category_code;
ALTER TABLE refunds ADD COLUMN merchant_category_code INT64;