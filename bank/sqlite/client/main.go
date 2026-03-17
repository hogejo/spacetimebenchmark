package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	config := parseConfig()
	accounts, initialBalance, requests, err := readRequests(config.input)
	if err != nil {
		log.Fatal(err)
	}
	if len(requests) == 0 {
		log.Fatal("no requests?!")
	}
	config.accounts = accounts
	config.initialBalance = initialBalance
	printSummary(&config)
	if !config.verifyOnly {
		runBenchmark(config, requests)
	}
	log.Println("Verifying the database ...")
	if err := runVerifications(config); err != nil {
		log.Fatal(err)
	}
}
