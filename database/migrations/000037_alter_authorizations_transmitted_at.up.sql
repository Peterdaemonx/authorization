ALTER TABLE authorizations DROP COLUMN transmitted_at;
ALTER TABLE authorizations ADD COLUMN transmitted_at TIMESTAMP;
ALTER TABLE authorizations ALTER COLUMN transmitted_at TIMESTAMP NOT NULL;