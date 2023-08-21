ALTER TABLE mastercard_refunds ADD COLUMN card_product_id STRING(3);
ALTER TABLE mastercard_refunds ADD COLUMN card_program_id STRING(3);
ALTER TABLE mastercard_refunds ALTER COLUMN card_product_id STRING(3) NOT NULL;
ALTER TABLE mastercard_refunds ALTER COLUMN card_program_id STRING(3) NOT NULL;

ALTER TABLE refunds DROP COLUMN card_product_id;
ALTER TABLE refunds DROP COLUMN card_program_id;