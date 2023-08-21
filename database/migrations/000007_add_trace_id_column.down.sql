ALTER TABLE authorizations DROP COLUMN initial_trace_id;
ALTER TABLE authorizations ADD COLUMN external_trace_id STRING(30);