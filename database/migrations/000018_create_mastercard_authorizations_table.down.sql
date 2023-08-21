ALTER TABLE authorizations DROP COLUMN response_code;
ALTER TABLE authorizations ADD COLUMN mc_system_trace_audit_number INT64;
ALTER TABLE authorizations ADD COLUMN mc_network_reporting_date STRING(100);
ALTER TABLE authorizations ADD COLUMN mc_response_code STRING(2);
ALTER TABLE authorizations ADD COLUMN financial_network_code STRING(100);
ALTER TABLE authorizations ADD COLUMN banknet_reference_number STRING(100);

DROP TABLE mastercard_authorizations;