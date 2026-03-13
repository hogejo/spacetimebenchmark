DROP TABLE IF EXISTS balances CASCADE;
CREATE TABLE balances
(
    id      BIGINT PRIMARY KEY,
    balance BIGINT
);
CREATE INDEX balances_id_index ON balances USING HASH (id);
