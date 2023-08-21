ALTER TABLE authorizations
ADD COLUMN retrieval_reference_number string(12);

ALTER TABLE visa_authorizations DROP COLUMN retrieval_reference_number;