INSERT INTO permissions (permission_id, code, label)
VALUES ("6451e565-bc86-4279-a11b-2c14c6893a1f", "create_authorization", "Create new authorization");
INSERT INTO permissions (permission_id, code, label)
VALUES ("af697e64-dbdb-442f-8542-582b32177719", "create_refund", "Create new refund");
INSERT INTO permissions (permission_id, code, label)
VALUES ("75dfed92-1006-4672-abe0-67deba0a0b55", "get_authorizations", "Get authorizations");
INSERT INTO permissions (permission_id, code, label)
VALUES ("b93ac092-b8fb-43cb-9a2d-fd099c8ddc05", "create_reversal", "Create reversal");
INSERT INTO permissions (permission_id, code, label)
VALUES ("be88a63e-dcd8-49f8-9dee-41725caa65de", "get_captures", "Get captures");
INSERT INTO permissions (permission_id, code, label)
VALUES ("44df31eb-3384-4f5a-9335-8207cbe9d45d", "get_refunds", "Get refunds");
INSERT INTO permissions (permission_id, code, label)
VALUES ("37a71765-99a0-44e0-b536-58cf0f973eec", "create_capture", "Create new capture");

INSERT psp (psp_id, name, prefix)
VALUES ("1779edcd-4f14-4c97-a61e-29a827e7ed89", "mycompany.com PS", "002");
INSERT psp (psp_id, name, prefix)
VALUES ("0cd8d732-66c2-4dae-bb99-16494dea7796", "test_psp", "001");

INSERT api_consumers (psp_id, api_consumer_id, name, api_key)
VALUES ("1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c","mycompany.com PS all permissions", "d035eb4f-0f35-4266-a837-84afe0a474a1");
INSERT api_consumers (psp_id, api_consumer_id, name, api_key)
VALUES ("0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d", "test_consumer_psp", "6247c10c-84a0-4fa1-b330-77eea1e944d3");

INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("6451e565-bc86-4279-a11b-2c14c6893a1f", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("af697e64-dbdb-442f-8542-582b32177719", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("37a71765-99a0-44e0-b536-58cf0f973eec", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("be88a63e-dcd8-49f8-9dee-41725caa65de", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("b93ac092-b8fb-43cb-9a2d-fd099c8ddc05", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("75dfed92-1006-4672-abe0-67deba0a0b55", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("44df31eb-3384-4f5a-9335-8207cbe9d45d", "1779edcd-4f14-4c97-a61e-29a827e7ed89", "7ed6c825-aef2-4a65-8b01-be25e778f48c");

INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("6451e565-bc86-4279-a11b-2c14c6893a1f", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("af697e64-dbdb-442f-8542-582b32177719", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("37a71765-99a0-44e0-b536-58cf0f973eec", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("be88a63e-dcd8-49f8-9dee-41725caa65de", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("b93ac092-b8fb-43cb-9a2d-fd099c8ddc05", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("44df31eb-3384-4f5a-9335-8207cbe9d45d", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");
INSERT INTO api_consumers_permissions (permission_id, psp_id, api_consumer_id)
VALUES ("75dfed92-1006-4672-abe0-67deba0a0b55", "0cd8d732-66c2-4dae-bb99-16494dea7796", "26c6c05d-b9f0-4e96-8265-e18df255197d");

INSERT INTO sequences (name, next_value, rollover_value)
VALUES ('mastercard_stan', 1000000, '20220627');
INSERT INTO sequences (name, next_value, rollover_value)
VALUES ('visa_stan', 1000000, '20220627');

INSERT INTO users (person_guid, username)
VALUES
    ('2c6984f3-af80-4bfd-81c7-0f597c12662e', 'test_account'),
    ('78b94843-5362-4760-96be-5c8b7df9c3fb', 'jorygeerts'),
    ('f005b430-049f-452d-a89e-262d7269853d', 'dannylucas');

INSERT INTO users_permissions (permission_id, person_guid)
SELECT p.permission_id, u.person_guid
FROM
    permissions AS p
        JOIN users AS u ON 1 = 1
 -- WHERE u.person_guid NOT IN (SELECT DISTINCT person_guid FROM users_permissions)
;


