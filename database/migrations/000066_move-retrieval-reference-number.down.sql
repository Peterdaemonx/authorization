ALTER TABLE authorizations DROP COLUMN retrieval_reference_number;

ALTER TABLE visa_authorizations ADD COLUMN retrieval_reference_number string(12);