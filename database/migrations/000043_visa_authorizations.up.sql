CREATE TABLE visa_authorizations
(
    authorization_id                              STRING(36) NOT NULL,
    acquiring_institution_country_code            STRING(3) NOT NULL,
    acquiring_institution_id                      INT64     NOT NULL,
    point_of_service_pan_entry_mode               STRING(2) NOT NULL,
    point_of_service_pin_entry_mode               STRING(2) NOT NULL,
    point_of_service_condition_code               STRING(2) NOT NULL,
    retrieval_reference_number                    STRING(12) NOT NULL,
    additional_pos_info_terminal_type             STRING(1) NOT NULL,
    additional_pos_info_terminal_entry_capability STRING(1) NOT NULL,
    additional_pos_info_type_or_level_indicator   STRING(2) NOT NULL,
    created_at                                    TIMESTAMP NOT NULL,
) PRIMARY KEY(authorization_id),
INTERLEAVE IN PARENT authorizations ON
DELETE
CASCADE;

ALTER TABLE mastercard_authorizations DROP COLUMN cardholder_transaction_type_code;
ALTER TABLE mastercard_authorizations DROP COLUMN cardholder_from_account_type_code;
ALTER TABLE mastercard_authorizations DROP COLUMN cardholder_to_account_type_code;

ALTER TABLE authorizations
    ADD COLUMN cardholder_transaction_type_code STRING(2);
ALTER TABLE authorizations
    ADD COLUMN cardholder_from_account_type_code STRING(2);
ALTER TABLE authorizations
    ADD COLUMN cardholder_to_account_type_code STRING(2);