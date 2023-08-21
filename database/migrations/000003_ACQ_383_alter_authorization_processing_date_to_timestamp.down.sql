ALTER TABLE authorizations DROP COLUMN transmitted_at;
ALTER TABLE authorizations ADD COLUMN transmitted_on DATE;