ALTER TABLE authorizations ADD COLUMN transaction_initiated_by STRING(17);
ALTER TABLE authorizations ADD COLUMN transaction_subcategory STRING(27);
