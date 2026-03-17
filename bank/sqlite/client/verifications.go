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
	var totalBalance uint64
	for i := uint64(0); i < config.accounts; i++ {
		response, err := conn.sendGet(i)
		if err != nil {
			return fmt.Errorf("get balance for account %d: %w", i, err)
		}
		balance, err := strconv.ParseUint(response, 10, 64)
		if err != nil {
			return fmt.Errorf("parse balance for account %d: %q: %w", i, response, err)
		}
		totalBalance += balance
	}
	expected := config.accounts * config.initialBalance
	if totalBalance != expected {
		return fmt.Errorf("balance mismatch: got %d, expected %d", totalBalance, expected)
	}
	fmt.Printf("Balance check passed (%d == %d expected).\n", totalBalance, expected)
	return nil
}

func verifyAccounts(conn *Connection, config Config) error {
	var totalAccounts uint64
	for i := uint64(0); i < config.accounts; i++ {
		response, err := conn.sendGet(i)
		if err != nil {
			return fmt.Errorf("get account %d: %w", i, err)
		}
		if _, err := strconv.ParseUint(response, 10, 64); err != nil {
			return fmt.Errorf("account %d returned unexpected response: %q", i, response)
		}
		totalAccounts++
	}
	expected := config.accounts
	if totalAccounts != expected {
		return fmt.Errorf("account number mismatch: got %d, expected %d", totalAccounts, expected)
	}
	fmt.Printf("Account check passed (%d == %d expected).\n", totalAccounts, expected)
	return nil
}
