ALTER TABLE authorizations DROP COLUMN mc_response_code;
ALTER TABLE authorizations ADD COLUMN mc_response_code STRING(2);