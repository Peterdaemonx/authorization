ALTER TABLE authorizations ALTER COLUMN merchant_name STRING(22) NOT NULL;
ALTER TABLE authorizations ALTER COLUMN merchant_city STRING(13) NOT NULL;
ALTER TABLE authorizations ALTER COLUMN merchant_country STRING(3) NOT NULL;
ALTER TABLE authorizations ALTER COLUMN merchant_category_code INT64 NOT NULL;
ALTER TABLE authorizations ALTER COLUMN merchant_card_acceptor_id STRING(12) NOT NULL;