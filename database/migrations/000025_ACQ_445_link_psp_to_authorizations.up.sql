ALTER TABLE authorizations ADD COLUMN psp_id STRING(36);
ALTER TABLE authorizations ALTER COLUMN psp_id STRING(36) NOT NULL;
ALTER TABLE authorizations ADD CONSTRAINT FK_psp_authorizations FOREIGN KEY (psp_id) REFERENCES psp (psp_id);