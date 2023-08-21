-- At the time this migration was created we didn't have a specific table for visa refunds.
-- so these migrations are created to keep in mind that the
-- point_of_service_pin_entry_mode for visa should be a length of 2

ALTER TABLE mastercard_refunds DROP COLUMN point_of_service_pan_entry_mode;
ALTER TABLE mastercard_refunds DROP COLUMN point_of_service_pin_entry_mode;

ALTER TABLE refunds ADD COLUMN point_of_service_pan_entry_mode STRING(2);
ALTER TABLE refunds ADD COLUMN point_of_service_pin_entry_mode STRING(2);
ALTER TABLE refunds ALTER COLUMN point_of_service_pan_entry_mode STRING(2) NOT NULL;
ALTER TABLE refunds ALTER COLUMN point_of_service_pin_entry_mode STRING(2) NOT NULL;
