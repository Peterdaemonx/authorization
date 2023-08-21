ALTER TABLE authorizations DROP COLUMN point_of_service_pan_entry_mode;
ALTER TABLE authorizations DROP COLUMN point_of_service_pin_entry_mode;

ALTER TABLE mastercard_authorizations ADD COLUMN point_of_service_pan_entry_mode STRING(2);
ALTER TABLE mastercard_authorizations ADD COLUMN point_of_service_pin_entry_mode STRING(1);
ALTER TABLE mastercard_authorizations ALTER COLUMN point_of_service_pin_entry_mode STRING(2) NOT NULL;
ALTER TABLE mastercard_authorizations ALTER COLUMN point_of_service_pin_entry_mode STRING(1) NOT NULL;

ALTER TABLE visa_authorizations ADD COLUMN point_of_service_pan_entry_mode STRING(2);
ALTER TABLE visa_authorizations ADD COLUMN point_of_service_pin_entry_mode STRING(2);
ALTER TABLE mastercard_authorizations ALTER COLUMN point_of_service_pin_entry_mode STRING(2) NOT NULL;
ALTER TABLE mastercard_authorizations ALTER COLUMN point_of_service_pin_entry_mode STRING(1) NOT NULL;

