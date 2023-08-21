ALTER TABLE authorization_reversals ADD COLUMN amount INT64;
ALTER TABLE authorization_reversals ALTER COLUMN amount INT64 NOT NULL;