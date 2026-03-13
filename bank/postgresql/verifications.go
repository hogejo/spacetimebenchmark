package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func runVerifications(pool *pgxpool.Pool, config Config) error {
	err := verifyTotalBalance(pool, config)
	if err != nil {
		return err
	}
	return verifyAccounts(pool, config)
}

func verifyTotalBalance(pool *pgxpool.Pool, config Config) error {
	var totalBalance uint64
	err := pool.QueryRow(context.Background(), "SELECT SUM(balance) FROM balances").Scan(&totalBalance)
	if err != nil {
		return fmt.Errorf("query total balance: %w", err)
	}
	expected := config.accounts * config.initialBalance
	if totalBalance != expected {
		return fmt.Errorf("balance mismatch: got %d, expected %d", totalBalance, expected)
	}
	fmt.Printf("Balance check passed (%d == %d expected).\n", totalBalance, expected)
	return nil
}

func verifyAccounts(pool *pgxpool.Pool, config Config) error {
	var totalAccounts uint64
	err := pool.QueryRow(context.Background(), "SELECT COUNT(id) FROM balances").Scan(&totalAccounts)
	if err != nil {
		return fmt.Errorf("query total accounts: %w", err)
	}
	expected := config.accounts
	if totalAccounts != expected {
		return fmt.Errorf("account number mismatch: got %d, expected %d", totalAccounts, expected)
	}
	fmt.Printf("Account check passed (%d == %d expected).\n", totalAccounts, expected)
	return nil
}
