ALTER TABLE authorizations DROP COLUMN update_at;
ALTER TABLE authorizations ADD COLUMN updated_at TIMESTAMP;