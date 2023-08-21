-- Create a table to store sequences in
CREATE TABLE sequences
(
    name           STRING(64) NOT NULL,
    next_value     INT64      NOT NULL,
    rollover_value STRING(64)
) PRIMARY KEY (name)
