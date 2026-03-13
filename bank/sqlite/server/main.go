package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(0)
	config := parseConfig()
	printSummary(&config)
	db, err := openDatabase(config.databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer closeDatabase(db)
	if err := prepareDatabase(db, config.accounts, config.initialBalance); err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runServer(ctx, config, db)
}
