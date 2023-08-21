ALTER TABLE authorizations DROP COLUMN mc_system_trace_audit_number;
ALTER TABLE authorizations DROP COLUMN mc_network_reporting_date;
ALTER TABLE authorizations DROP COLUMN mc_response_code;
ALTER TABLE authorizations DROP COLUMN financial_network_code;
ALTER TABLE authorizations DROP COLUMN banknet_reference_number;
ALTER TABLE authorizations ADD COLUMN response_code STRING(2);
ALTER TABLE authorizations ADD COLUMN system_trace_audit_number INT64;

CREATE TABLE mastercard_authorizations
(
    authorization_id STRING(36) NOT NULL,
    cardholder_transaction_type_code STRING(2) NOT NULL,
    cardholder_from_account_type_code STRING(2) NOT NULL,
    cardholder_to_account_type_code STRING(2) NOT NULL,
    point_of_service_pan_entry_mode STRING(2) NOT NULL,
    point_of_service_pin_entry_mode STRING(1) NOT NULL,
    pos_pin_capture_code STRING(2),
    pin_service_code STRING(2),
    terminal_attendance STRING(1) NOT NULL,
    terminal_location INT64 NOT NULL,
    cardholder_presence INT64 NOT NULL,
    card_presence INT64 NOT NULL,
    card_capture_capabilities INT64 NOT NULL,
    cardholder_activated_terminal_level INT64 NOT NULL,
    card_data_terminal_input_capability_indicator INT64,
    security_protocol INT64 NOT NULL,
    cardholder_authentication INT64 NOT NULL,
    ucaf_collection_indicator INT64 NOT NULL,
    network_reporting_date STRING(100),
    financial_network_code STRING(100),
    banknet_reference_number STRING(100),
    created_at TIMESTAMP NOT NULL,
) PRIMARY KEY(authorization_id),
INTERLEAVE IN PARENT authorizations ON DELETE CASCADE;