CREATE OR REPLACE PROCEDURE transfer(from_id BIGINT, to_id BIGINT, amount BIGINT)
LANGUAGE plpgsql
AS $$
DECLARE
    from_balance   BIGINT;
BEGIN
    IF from_id = to_id THEN
        RAISE EXCEPTION 'same_account';
    END IF;

    IF amount <= 0 THEN
        RAISE EXCEPTION 'invalid_amount';
    END IF;

    -- Lock both rows in a deterministic order (the lowest id first)
    -- to prevent deadlocks across concurrent calls.
    PERFORM 1 FROM balances WHERE id = LEAST(from_id, to_id) FOR UPDATE;
    PERFORM 1 FROM balances WHERE id = GREATEST(from_id, to_id) FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'account_not_found';
    END IF;

    SELECT balance
    INTO STRICT from_balance
    FROM balances
    WHERE id = from_id;
    IF from_balance < amount THEN
        RAISE EXCEPTION 'insufficient_funds';
    END IF;

    -- Apply the transfer in a single statement
    UPDATE balances
    SET balance = CASE id
      WHEN from_id THEN balance - amount
      WHEN to_id   THEN balance + amount
    END
    WHERE id IN (from_id, to_id);
EXCEPTION
    WHEN no_data_found THEN
        RAISE EXCEPTION 'account_not_found';
END;
$$;