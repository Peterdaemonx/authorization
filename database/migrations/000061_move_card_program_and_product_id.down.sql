ALTER TABLE refunds ADD COLUMN card_product_id STRING(3);
ALTER TABLE refunds ADD COLUMN card_program_id STRING(3);

ALTER TABLE mastercard_refunds DROP COLUMN card_product_id;
ALTER TABLE mastercard_refunds DROP COLUMN card_program_id;