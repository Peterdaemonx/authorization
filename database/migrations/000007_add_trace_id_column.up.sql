ALTER TABLE authorizations DROP COLUMN external_trace_id;
ALTER TABLE authorizations ADD COLUMN initial_trace_id STRING(30);