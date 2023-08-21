ALTER TABLE visa_authorizations ADD COLUMN chip_condition_code STRING(1);
ALTER TABLE visa_authorizations ADD COLUMN special_condition_indicator STRING(1);
ALTER TABLE visa_authorizations ADD COLUMN chip_transaction_indicator STRING(1);
ALTER TABLE visa_authorizations ADD COLUMN chip_card_authentication_reliability_indicator STRING(1);
ALTER TABLE visa_authorizations ADD COLUMN cardholder_id_method_indicator STRING(1);
ALTER TABLE visa_authorizations ADD COLUMN additional_authorization_indicators STRING(1);