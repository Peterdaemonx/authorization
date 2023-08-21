ALTER TABLE authorizations DROP COLUMN merchant_category_code;
ALTER TABLE authorizations ADD COLUMN merchant_category_code string(4);
ALTER TABLE authorizations ALTER COLUMN merchant_category_code string(4) NOT NULL;

ALTER TABLE refunds DROP COLUMN merchant_category_code;
ALTER TABLE refunds ADD COLUMN merchant_category_code string(4);
ALTER TABLE refunds ALTER COLUMN merchant_category_code string(4) NOT NULL;
