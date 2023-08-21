drop table authorization_reversals;

drop table authorizations;

CREATE TABLE authorizations (
    id STRING(36) NOT NULL,
    status STRING(50) NOT NULL,
    merchant_id STRING(36) NOT NULL,
    masked_pan STRING(40) NOT NULL,
    pan_token_id STRING(40) NOT NULL,
    card_scheme STRING(30) NOT NULL,
    amount INT64 NOT NULL,
    currency STRING(3) NOT NULL,
    localdatetime TIMESTAMP NOT NULL,
    country STRING(3),
    cardholder_verification STRING(30) NOT NULL,
    `type` STRING(30),
    `source` STRING(20) NOT NULL,
    exemption_type STRING(25),
    exemption_reason STRING(15),
    is_initial_recurring BOOL,
    threeds_version STRING(5),
    threeds_original_ecommerce_indicator STRING(3),
    threeds_applied_ecommerce_indicator STRING(3),
    reason_downgrade STRING(1),
    financial_network_code STRING(100),
    banknet_reference_number STRING(100),
    mc_system_trace_audit_number INT64,
    mc_network_reporting_date STRING(100),
    created_at TIMESTAMP NOT NULL,
    transmitted_at TIMESTAMP,
    customer_reference STRING(100),
    mc_response_code STRING(2),
    updated_at TIMESTAMP,
    initial_trace_id STRING(30),
) PRIMARY KEY (id);

CREATE TABLE authorization_reversals
(
    id              STRING(36) NOT NULL,
    reversal_id     STRING(36) NOT NULL,
    status          STRING(50) NOT NULL,
    stan            INT64,
    response_code   STRING(2),
    created_at      TIMESTAMP NOT NULL,
) PRIMARY KEY(id, reversal_id),
  INTERLEAVE IN PARENT authorizations ON DELETE CASCADE;