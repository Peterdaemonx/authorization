ALTER TABLE authorizations DROP COLUMN transmitted_on;
ALTER TABLE authorizations ADD COLUMN transmitted_at TIMESTAMP;
