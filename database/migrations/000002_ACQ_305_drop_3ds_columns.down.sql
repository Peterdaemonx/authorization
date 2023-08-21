ALTER TABLE authorizations ADD COLUMN threeds_authentication_verification_value STRING(36);
ALTER TABLE authorizations ADD COLUMN threeds_directory_server_transaction_id STRING(100);