ALTER TABLE authorizations DROP COLUMN threeds_original_ecommerce_indicator;
ALTER TABLE authorizations DROP COLUMN threeds_applied_ecommerce_indicator;

ALTER TABLE authorizations ADD COLUMN threeds_original_ecommerce_indicator STRING(3);
ALTER TABLE authorizations ADD COLUMN threeds_applied_ecommerce_indicator STRING(3);