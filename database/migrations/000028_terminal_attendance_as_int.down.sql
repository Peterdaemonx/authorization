ALTER TABLE mastercard_authorizations DROP COLUMN terminal_attendance;
ALTER TABLE mastercard_authorizations ADD COLUMN terminal_attendance STRING(1);