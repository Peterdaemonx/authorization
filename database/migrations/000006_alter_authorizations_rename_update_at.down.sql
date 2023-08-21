ALTER TABLE authorizations DROP COLUMN updated_at;
ALTER TABLE authorizations ADD COLUMN update_at TIMESTAMP;
