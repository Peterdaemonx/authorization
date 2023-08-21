ALTER TABLE authorizations ADD COLUMN card_sequence STRING(3);
ALTER TABLE authorizations ADD COLUMN terminal_id STRING(8);
ALTER TABLE authorizations ADD COLUMN terminal_capability STRING(50);
ALTER TABLE authorizations ADD COLUMN card_holder_verification_method STRING(15);

ALTER TABLE authorizations ADD COLUMN card_holder_activated_terminal_level STRING(50);