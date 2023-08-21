ALTER TABLE mastercard_authorizations ADD COLUMN authorization_id_response STRING(6);

ALTER TABLE authorizations DROP COLUMN authorization_id_response;