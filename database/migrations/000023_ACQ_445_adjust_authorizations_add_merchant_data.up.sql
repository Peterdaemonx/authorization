ALTER TABLE authorizations ADD COLUMN merchant_name STRING(22);
ALTER TABLE authorizations ALTER COLUMN merchant_name STRING(22) NOT NULL;
ALTER TABLE authorizations ADD COLUMN merchant_city STRING(13);
ALTER TABLE authorizations ALTER COLUMN merchant_city STRING(13) NOT NULL;
ALTER TABLE authorizations ADD COLUMN merchant_country STRING(3);
ALTER TABLE authorizations ALTER COLUMN merchant_country STRING(3) NOT NULL;
ALTER TABLE authorizations ADD COLUMN merchant_category_code INT64;
ALTER TABLE authorizations ALTER COLUMN merchant_category_code INT64 NOT NULL;
ALTER TABLE authorizations ADD COLUMN merchant_card_acceptor_id STRING(12);
ALTER TABLE authorizations ALTER COLUMN merchant_card_acceptor_id STRING(12) NOT NULL;