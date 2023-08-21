ALTER TABLE authorizations DROP COLUMN transmitted_at;
ALTER TABLE authorizations ADD COLUMN transmitted_at DATE;