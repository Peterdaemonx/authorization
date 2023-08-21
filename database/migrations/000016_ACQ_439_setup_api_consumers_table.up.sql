CREATE TABLE api_consumers(
    psp_id STRING(36) NOT NULL,
    api_consumer_id STRING(36) NOT NULL,
    name STRING(100) NOT NULL,
    api_key STRING(255) NOT NULL,
) PRIMARY KEY (psp_id, api_consumer_id),
INTERLEAVE IN PARENT psp ON DELETE CASCADE;