ALTER TABLE mastercard_authorizations DROP COLUMN point_of_service_pan_entry_mode;
ALTER TABLE mastercard_authorizations DROP COLUMN point_of_service_pin_entry_mode;

ALTER TABLE visa_authorizations DROP COLUMN point_of_service_pan_entry_mode;
ALTER TABLE visa_authorizations DROP COLUMN point_of_service_pin_entry_mode;

ALTER TABLE authorizations ADD COLUMN point_of_service_pan_entry_mode STRING(2);
ALTER TABLE authorizations ADD COLUMN point_of_service_pin_entry_mode STRING(2);

ALTER TABLE authorizations ALTER COLUMN point_of_service_pan_entry_mode STRING(2) NOT NULL;
ALTER TABLE authorizations ALTER COLUMN point_of_service_pin_entry_mode STRING(2) NOT NULL;
