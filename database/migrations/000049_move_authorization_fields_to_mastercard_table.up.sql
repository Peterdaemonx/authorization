ALTER TABLE authorizations DROP COLUMN card_product_id;
ALTER TABLE authorizations DROP COLUMN card_program_id;

ALTER TABLE mastercard_authorizations ADD COLUMN card_product_id STRING(3);
ALTER TABLE mastercard_authorizations ALTER COLUMN card_product_id STRING(3) NOT NULL;
ALTER TABLE mastercard_authorizations ADD COLUMN card_program_id STRING(3);
ALTER TABLE mastercard_authorizations ALTER COLUMN card_program_id STRING(3) NOT NULL;
