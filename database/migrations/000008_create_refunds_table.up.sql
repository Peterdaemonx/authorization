CREATE TABLE refunds
(
    id                                        STRING(36) NOT NULL,
    status                                    STRING(50) NOT NULL,
    merchant_id                               STRING(36) NOT NULL,
    masked_pan                                STRING(40) NOT NULL,
    pan_token_id                              STRING(40) NOT NULL,
    card_scheme                               STRING(30) NOT NULL,
    amount                                    INT64     NOT NULL,
    currency                                  STRING(3) NOT NULL,
    localdatetime                             TIMESTAMP NOT NULL,
    country                                   STRING(3),
    cardholder_verification                   STRING(30) NOT NULL,
    source                                    STRING(20) NOT NULL,
    financial_network_code                    STRING(100),
    banknet_reference_number                  STRING(100),
    mc_response_code                          INT64,
    mc_system_trace_audit_number              INT64,
    mc_network_reporting_date                 STRING(100),
    transmitted_on                            DATE,
    created_at                                TIMESTAMP NOT NULL,
    updated_at                                TIMESTAMP
) PRIMARY KEY(id);