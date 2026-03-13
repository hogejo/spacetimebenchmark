package main

import (
	"flag"
	"log"
)

const (
	defaultAddr           = ":9000"
	defaultDatabaseFile   = "bank.sqlite"
	defaultAccounts       = 1_000_000
	defaultInitialBalance = 100_000
)

type Config struct {
	addr           string
	databaseFile   string
	accounts       uint64
	initialBalance uint64
}

func parseConfig() Config {
	config := Config{}
	flag.StringVar(&config.addr, "addr", defaultAddr, "listen address")
	flag.StringVar(&config.databaseFile, "database-file", defaultDatabaseFile, "path to SQLite database file")
	flag.Uint64Var(&config.accounts, "accounts", defaultAccounts, "number of accounts")
	flag.Uint64Var(&config.initialBalance, "initial-balance", defaultInitialBalance, "initial balance per account")
	flag.Parse()
	return config
}

func printSummary(config *Config) {
	log.Printf("Server configuration summary:\n")
	log.Printf("  listen address:  %s\n", config.addr)
	log.Printf("  database file:   %s\n", config.databaseFile)
	log.Printf("  accounts:        %d\n", config.accounts)
	log.Printf("  initial balance: %d\n", config.initialBalance)
}
