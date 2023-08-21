CREATE TABLE authorizations
(
        id STRING(36) NOT NULL,
        status STRING(50) NOT NULL,
        merchant_id STRING(36) NOT NULL,
        masked_pan STRING(40) NOT NULL,
        pan_token_id STRING(40) NOT NULL,
        card_scheme STRING(30) NOT NULL,
        amount INT64 NOT NULL,
        currency STRING(3) NOT NULL,
        localdatetime TIMESTAMP NOT NULL,
        external_trace_id STRING(15),
        country STRING(3),
        cardholder_verification STRING(30) NOT NULL,
        type STRING(30),
        source STRING(20) NOT NULL,
        exemption_type STRING(25),
        exemption_reason STRING(15),
        is_initial_recurring BOOL,
        threeds_version STRING(5),
        threeds_authentication_verification_value STRING(36),
        threeds_directory_server_transaction_id STRING(100),
        threeds_original_ecommerce_indicator STRING(3),
        threeds_applied_ecommerce_indicator STRING(3),
        reason_downgrade STRING(1),
        financial_network_code STRING(100),
        banknet_reference_number STRING(100),
        mc_response_code INT64,
        mc_system_trace_audit_number INT64,
        mc_network_reporting_date STRING(100),
        transmitted_on DATE,
        created_at TIMESTAMP NOT NULL,
        update_at TIMESTAMP,
) PRIMARY KEY(id);