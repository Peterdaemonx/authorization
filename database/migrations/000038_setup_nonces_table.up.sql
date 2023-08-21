CREATE TABLE nonces (
    psp_id STRING(36) NOT NULL,
    nonce STRING(100) NOT NULL,
    used_at TIMESTAMP NOT NULL,
    CONSTRAINT FK_psp_nonces FOREIGN KEY (psp_id) REFERENCES psp (psp_id)
) PRIMARY KEY(psp_id, nonce);