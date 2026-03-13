UPDATE balances
SET balance = CASE id
    WHEN ?1 THEN balance - ?3
    WHEN ?2 THEN balance + ?3
END
WHERE id IN (?1, ?2);