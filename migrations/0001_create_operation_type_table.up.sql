CREATE TABLE operation_types (
    id SMALLINT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
INSERT INTO operation_types (id, name) VALUES
(1, 'DEPOSIT'),
(2, 'WITHDRAW');
