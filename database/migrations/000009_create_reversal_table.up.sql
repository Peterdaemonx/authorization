CREATE TABLE authorization_reversals
(
    id                                        STRING(36) NOT NULL,
    status                                    STRING(50) NOT NULL,
    stan                                      INT64,
    authorization_id                          STRING(36) NOT NULL,
    response_code                             STRING(2),
    created_at                                TIMESTAMP NOT NULL,
) PRIMARY KEY(id);