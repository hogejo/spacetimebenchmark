DROP TABLE IF EXISTS balances;
CREATE TABLE balances
(
    id      INTEGER PRIMARY KEY,
    balance INTEGER
);
CREATE INDEX balances_id_index ON balances (id);
