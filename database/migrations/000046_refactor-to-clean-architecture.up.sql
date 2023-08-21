ALTER TABLE authorizations DROP COLUMN country;
ALTER TABLE authorizations DROP COLUMN merchant_name;
ALTER TABLE authorizations ADD COLUMN card_acceptor_name STRING(22);
ALTER TABLE authorizations ALTER COLUMN card_acceptor_name STRING(22) NOT NULL;
ALTER TABLE authorizations DROP COLUMN merchant_city;
ALTER TABLE authorizations ADD COLUMN card_acceptor_city STRING(13);
ALTER TABLE authorizations ALTER COLUMN card_acceptor_city STRING(13) NOT NULL;
ALTER TABLE authorizations DROP COLUMN merchant_country;
ALTER TABLE authorizations ADD COLUMN card_acceptor_country STRING(3);
ALTER TABLE authorizations ALTER COLUMN card_acceptor_country STRING(3) NOT NULL;
ALTER TABLE authorizations DROP COLUMN merchant_card_acceptor_id;
ALTER TABLE authorizations ADD COLUMN card_acceptor_id STRING(12);
ALTER TABLE authorizations ALTER COLUMN card_acceptor_id STRING(12) NOT NULL;
ALTER TABLE authorizations ADD COLUMN card_acceptor_postal_code STRING(10);
ALTER TABLE authorizations DROP COLUMN merchant_postal_code;

ALTER TABLE authorizations ADD COLUMN card_acceptor_category_code STRING(4);
ALTER TABLE authorizations ALTER COLUMN card_acceptor_category_code STRING(4) NOT NULL;
ALTER TABLE authorizations DROP COLUMN merchant_category_code;

ALTER TABLE authorizations DROP COLUMN type;
ALTER TABLE mastercard_authorizations ADD COLUMN authorization_type STRING(30);
ALTER TABLE mastercard_authorizations ALTER COLUMN authorization_type STRING(30) NOT NULL;
ALTER TABLE authorizations DROP COLUMN exemption_type;
ALTER TABLE authorizations ADD COLUMN exemption STRING(25);
ALTER TABLE authorizations DROP COLUMN exemption_reason;

ALTER TABLE authorizations ALTER COLUMN transmitted_at TIMESTAMP;

ALTER TABLE authorizations DROP COLUMN threeds_authentication_verification_value;
ALTER TABLE authorizations DROP COLUMN trace_id;

ALTER TABLE mastercard_authorizations ADD COLUMN authorization_id_response STRING(6);

ALTER TABLE authorization_captures ADD COLUMN reference STRING(100);

ALTER TABLE mastercard_authorizations DROP COLUMN terminal_attendance;
ALTER TABLE mastercard_authorizations DROP COLUMN terminal_location;
ALTER TABLE mastercard_authorizations DROP COLUMN cardholder_presence;
ALTER TABLE mastercard_authorizations DROP COLUMN card_presence;
ALTER TABLE mastercard_authorizations DROP COLUMN card_capture_capabilities;
ALTER TABLE mastercard_authorizations DROP COLUMN cardholder_activated_terminal_level;
ALTER TABLE mastercard_authorizations DROP COLUMN card_data_terminal_input_capability_indicator;

ALTER TABLE mastercard_authorizations ADD COLUMN terminal_attendance INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN terminal_location INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_holder_presence INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_presence INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_capture_capabilities INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN transaction_status INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN transaction_security INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_holder_activated_terminal_level INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN card_data_terminal_input_capability_indicator INT64;
ALTER TABLE mastercard_authorizations ADD COLUMN authorization_life_cycle STRING(2);
ALTER TABLE mastercard_authorizations ADD COLUMN country_code STRING(3);
ALTER TABLE mastercard_authorizations ADD COLUMN postal_code STRING(10);

DROP TABLE refunds;

CREATE TABLE refunds (
    refund_id                    STRING(36) NOT NULL,
    status                       STRING(50) NOT NULL,
    masked_pan                   STRING(40) NOT NULL,
    pan_token_id                 STRING(40) NOT NULL,
    card_scheme                  STRING(30) NOT NULL,
    amount                       INT64 NOT NULL,
    currency                     STRING(3) NOT NULL,
    localdatetime                TIMESTAMP NOT NULL,
    source                       STRING(20) NOT NULL,
    system_trace_audit_number INT64,
    created_at                   TIMESTAMP NOT NULL,
    updated_at                   TIMESTAMP,
    response_code             STRING(20),
    psp_id                       STRING(36) NOT NULL,
    transmitted_at               TIMESTAMP,
    customer_reference           STRING(100),
    card_acceptor_name           STRING(22) NOT NULL,
    card_acceptor_city           STRING(13) NOT NULL,
    card_acceptor_country        STRING(3) NOT NULL,
    card_acceptor_id             STRING(12) NOT NULL,
    card_acceptor_postal_code    STRING(10),
    card_acceptor_category_code  STRING(4) NOT NULL,
    card_issuer_id               STRING(11),
    card_issuer_name             STRING(70),
    card_issuer_countrycode      STRING(3),
    card_product_id              STRING(3),
    card_program_id              STRING(3),
    cardholder_transaction_type_code        STRING(2),
        cardholder_from_account_type_code        STRING(2),
        cardholder_to_account_type_code        STRING(2)
) PRIMARY KEY (refund_id);

ALTER TABLE refunds ADD CONSTRAINT FK_psp_refunds FOREIGN KEY (psp_id) REFERENCES psp (psp_id);

CREATE TABLE mastercard_refunds (
    refund_id STRING(36) NOT NULL,
    point_of_service_pan_entry_mode               STRING(2) NOT NULL,
    point_of_service_pin_entry_mode               STRING(1) NOT NULL,
    network_reporting_date                        STRING(100),
    financial_network_code                        STRING(100),
    banknet_reference_number                      STRING(100),
    created_at                                    TIMESTAMP NOT NULL,
    reference                                     STRING(30),
    authorization_type                            STRING(30) NOT NULL,
    authorization_id_response                     STRING(6),
    terminal_attendance                           INT64,
    terminal_location                             INT64,
    card_holder_presence                          INT64,
    card_presence                                 INT64,
    card_capture_capabilities                     INT64,
    transaction_status                            INT64,
    transaction_security                          INT64,
    card_holder_activated_terminal_level          INT64,
    card_data_terminal_input_capability_indicator INT64,
    authorization_life_cycle                      STRING(2),
    country_code                                  STRING(3),
    postal_code                                   STRING(10)
) PRIMARY KEY(refund_id),
INTERLEAVE IN PARENT refunds ON DELETE CASCADE;

ALTER TABLE authorizations DROP COLUMN reason_downgrade;
ALTER TABLE authorizations DROP COLUMN threeds_applied_ecommerce_indicator;
ALTER TABLE mastercard_authorizations ADD COLUMN reason_ucaf_downgrade INT64;

ALTER TABLE visa_authorizations DROP COLUMN acquiring_institution_country_code;
ALTER TABLE visa_authorizations DROP COLUMN acquiring_institution_id;