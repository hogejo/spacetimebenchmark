package main

import (
	"fmt"
	"log"
	"strconv"
)

func runVerifications(config Config) error {
	conn, err := newConnection(config.serverAddr)
	if err != nil {
		return fmt.Errorf("connect for verification: %w", err)
	}
	defer conn.Close()

	log.Println("Verifying total balance ...")
	err = verifyTotalBalance(conn, config)
	if err != nil {
		return err
	}
	log.Println("Verifying numbers of accounts ...")
	return verifyAccounts(conn, config)
}

func verifyTotalBalance(conn *Connection, config Config) error {
	response, err := conn.sendGetTotal()
	if err != nil {
		return fmt.Errorf("get total balance: %w", err)
	}
	totalBalance, err := strconv.ParseUint(response, 10, 64)
	if err != nil {
		return fmt.Errorf("parse total balance: %q: %w", response, err)
	}
	expected := config.accounts * config.initialBalance
	if totalBalance != expected {
		return fmt.Errorf("balance mismatch: got %d, expected %d", totalBalance, expected)
	}
	fmt.Printf("Balance check passed (%d == %d expected).\n", totalBalance, expected)
	return nil
}

func verifyAccounts(conn *Connection, config Config) error {
	response, err := conn.sendCountAccounts(0, config.accounts-1)
	if err != nil {
		return fmt.Errorf("count accounts: %w", err)
	}
	totalAccounts, err := strconv.ParseUint(response, 10, 64)
	if err != nil {
		return fmt.Errorf("parse count_accounts: %q: %w", response, err)
	}
	expected := config.accounts
	if totalAccounts != expected {
		return fmt.Errorf("account number mismatch: got %d, expected %d", totalAccounts, expected)
	}
	fmt.Printf("Account check passed (%d == %d expected).\n", totalAccounts, expected)
	return nil
}
