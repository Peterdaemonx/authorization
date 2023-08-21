DROP TABLE visa_authorizations;

ALTER TABLE authorizations DROP COLUMN cardholder_transaction_type_code;
ALTER TABLE authorizations DROP COLUMN cardholder_from_account_type_code;
ALTER TABLE authorizations DROP COLUMN cardholder_to_account_type_code;

ALTER TABLE mastercard_authorizations ADD COLUMN cardholder_transaction_type_code STRING(2);
ALTER TABLE mastercard_authorizations ADD COLUMN cardholder_from_account_type_code STRING(2);
ALTER TABLE mastercard_authorizations ADD COLUMN cardholder_to_account_type_code STRING(2);