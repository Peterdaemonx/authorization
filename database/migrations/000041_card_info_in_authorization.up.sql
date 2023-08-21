
ALTER TABLE authorizations ADD COLUMN card_issuer_id STRING(10);
ALTER TABLE authorizations ADD COLUMN card_issuer_name STRING(70);
ALTER TABLE authorizations ADD COLUMN card_issuer_countrycode STRING(3);
ALTER TABLE authorizations ADD COLUMN card_product_id STRING(3);
ALTER TABLE authorizations ADD COLUMN card_program_id STRING(3);