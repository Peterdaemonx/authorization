ALTER TABLE authorizations ADD COLUMN country STRING(3);

ALTER TABLE authorizations ADD COLUMN merchant_name STRING(22);
ALTER TABLE authorizations ALTER COLUMN merchant_name STRING(22) NOT NULL;
ALTER TABLE authorizations DROP COLUMN card_acceptor_name;
ALTER TABLE authorizations ADD COLUMN merchant_city STRING(13);
ALTER TABLE authorizations ALTER COLUMN merchant_city STRING(13) NOT NULL;
ALTER TABLE authorizations DROP COLUMN card_acceptor_city;
ALTER TABLE authorizations ADD COLUMN merchant_country STRING(3);
ALTER TABLE authorizations ALTER COLUMN merchant_country STRING(3) NOT NULL;
ALTER TABLE authorizations DROP COLUMN card_acceptor_country;
ALTER TABLE authorizations ADD COLUMN merchant_card_acceptor_id STRING(12);
ALTER TABLE authorizations ALTER COLUMN merchant_card_acceptor_id STRING(12) NOT NULL;
ALTER TABLE authorizations DROP COLUMN card_acceptor_id;
ALTER TABLE authorizations ADD COLUMN merchant_postal_code STRING(10);
ALTER TABLE authorizations DROP COLUMN card_acceptor_postal_code;
ALTER TABLE authorizations ADD COLUMN merchant_category_code STRING(4);
ALTER TABLE authorizations ALTER COLUMN merchant_category_code STRING(4) NOT NULL;
ALTER TABLE authorizations DROP COLUMN card_acceptor_category_code;
ALTER TABLE authorizations ADD COLUMN threeds_applied_ecommerce_indicator INT64;

ALTER TABLE authorizations ADD COLUMN type STRING(30);
ALTER TABLE mastercard_authorizations DROP COLUMN authorization_type;

ALTER TABLE authorizations DROP COLUMN exemption;
ALTER TABLE authorizations ADD COLUMN exemption_type STRING(25);
ALTER TABLE authorizations ADD COLUMN exemption_reason STRING(15);

ALTER TABLE authorizations DROP COLUMN transmitted_at;
ALTER TABLE authorizations ADD COLUMN transmitted_at TIMESTAMP;
ALTER TABLE authorizations ALTER COLUMN transmitted_at TIMESTAMP NOT NULL;

ALTER TABLE authorizations ADD COLUMN threeds_authentication_verification_value STRING(36);
ALTER TABLE authorizations ADD COLUMN trace_id STRING(30);

ALTER TABLE mastercard_authorizations DROP COLUMN authorization_id_response;

ALTER TABLE authorization_captures DROP COLUMN reference;

ALTER TABLE mastercard_authorizations DROP COLUMN terminal_attendance;
ALTER TABLE mastercard_authorizations DROP COLUMN terminal_location;
ALTER TABLE mastercard_authorizations DROP COLUMN card_holder_presence;
ALTER TABLE mastercard_authorizations DROP COLUMN card_presence;
ALTER TABLE mastercard_authorizations DROP COLUMN card_capture_capabilities;
ALTER TABLE mastercard_authorizations DROP COLUMN transaction_status;
ALTER TABLE mastercard_authorizations DROP COLUMN transaction_security;
ALTER TABLE mastercard_authorizations DROP COLUMN card_holder_activated_terminal_level;
ALTER TABLE mastercard_authorizations DROP COLUMN card_data_terminal_input_capability_indicator;
ALTER TABLE mastercard_authorizations DROP COLUMN authorization_life_cycle;
ALTER TABLE mastercard_authorizations DROP COLUMN country_code;
ALTER TABLE mastercard_authorizations DROP COLUMN postal_code;

ALTER TABLE mastercard_authorizations ADD COLUMN terminal_attendance INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN terminal_location INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN cardholder_presence INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_presence INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_capture_capabilities INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN cardholder_activated_terminal_level INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_data_terminal_input_capability_indicator INT64;

ALTER TABLE refunds DROP COLUMN customer_reference;

ALTER TABLE refunds ADD COLUMN merchant_name STRING(22);
ALTER TABLE refunds ALTER COLUMN merchant_name STRING(22) NOT NULL;
ALTER TABLE refunds DROP COLUMN card_acceptor_name;
ALTER TABLE refunds ADD COLUMN merchant_city STRING(13);
ALTER TABLE refunds ALTER COLUMN merchant_city STRING(13) NOT NULL;
ALTER TABLE refunds DROP COLUMN card_acceptor_city;
ALTER TABLE refunds ADD COLUMN merchant_country STRING(3);
ALTER TABLE refunds ALTER COLUMN merchant_country STRING(3) NOT NULL;
ALTER TABLE refunds DROP COLUMN card_acceptor_country;
ALTER TABLE refunds ADD COLUMN merchant_card_acceptor_id STRING(12);
ALTER TABLE refunds ALTER COLUMN merchant_card_acceptor_id STRING(12) NOT NULL;
ALTER TABLE refunds DROP COLUMN card_acceptor_id;
ALTER TABLE refunds ADD COLUMN merchant_postal_code STRING(10);
ALTER TABLE refunds DROP COLUMN card_acceptor_postal_code;
ALTER TABLE refunds ADD COLUMN merchant_category_code STRING(4);
ALTER TABLE refunds ALTER COLUMN merchant_category_code STRING(4) NOT NULL;
ALTER TABLE refunds DROP COLUMN card_acceptor_category_code;
ALTER TABLE refunds ADD COLUMN mc_response_code INT64;

ALTER TABLE refunds DROP COLUMN card_issuer_id;
ALTER TABLE refunds DROP COLUMN card_issuer_name;
ALTER TABLE refunds DROP COLUMN card_issuer_countrycode;
ALTER TABLE refunds DROP COLUMN card_product_id;
ALTER TABLE refunds DROP COLUMN card_program_id;

DROP TABLE mastercard_refunds;