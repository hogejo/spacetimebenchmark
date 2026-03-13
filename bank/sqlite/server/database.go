package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/schema.sql
var schemaSQL string

//go:embed sql/transfer.sql
var transferSQL string

func openDatabase(databaseFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}
	log.Printf("opened sqlite database: %s", databaseFile)
	return db, nil
}

func closeDatabase(db *sql.DB) {
	db.Close()
}

func transfer(db *sql.DB, fromID, toID, amount uint64) error {
	if fromID == toID {
		return fmt.Errorf("same_account")
	}
	if amount == 0 {
		return fmt.Errorf("invalid_amount")
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock rows in deterministic order by reading the lowest id first
	var balance uint64
	lowID, highID := fromID, toID
	if lowID > highID {
		lowID, highID = highID, lowID
	}
	err = tx.QueryRow("SELECT balance FROM balances WHERE id = ?", lowID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("account_not_found")
	}
	var balance2 uint64
	err = tx.QueryRow("SELECT balance FROM balances WHERE id = ?", highID).Scan(&balance2)
	if err != nil {
		return fmt.Errorf("account_not_found")
	}

	// Determine from_balance
	fromBalance := balance
	if lowID != fromID {
		fromBalance = balance2
	}
	if fromBalance < amount {
		return fmt.Errorf("insufficient_funds")
	}

	result, err := tx.Exec(transferSQL, fromID, toID, amount)
	if err != nil {
		return fmt.Errorf("execute transfer: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows != 2 {
		return fmt.Errorf("account_not_found")
	}
	return tx.Commit()
}

func getBalance(db *sql.DB, accountID uint64) (uint64, error) {
	var balance uint64
	err := db.QueryRow("SELECT balance FROM balances WHERE id = ?", accountID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("account_not_found")
	}
	return balance, nil
}

func prepareDatabase(db *sql.DB, accounts, initialBalance uint64) error {
	fmt.Printf("preparing the database ...")
	if _, err := db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("prepare database: %w", err)
	}
	if _, err := db.Exec("DELETE FROM balances"); err != nil {
		return fmt.Errorf("clear balances: %w", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO balances (id, balance) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("prepare seed statement: %w", err)
	}
	defer stmt.Close()
	for i := uint64(0); i < accounts; i++ {
		if _, err := stmt.Exec(i, initialBalance); err != nil {
			return fmt.Errorf("seed balance %d: %w", i, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}
	fmt.Println(" done.")
	return nil
}
