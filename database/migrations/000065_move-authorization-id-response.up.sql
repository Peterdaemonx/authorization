ALTER TABLE authorizations
ADD COLUMN authorization_id_response string(6);

ALTER TABLE mastercard_authorizations DROP COLUMN authorization_id_response;