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
	pool, err := openPool(config)
	if err != nil {
		log.Fatal(err)
	}
	if !config.verifyOnly {
		if err := prepareDatabase(pool, config); err != nil {
			log.Fatal(err)
		}
		runBenchmark(pool, config, requests)
	}
	if err := runVerifications(pool, config); err != nil {
		log.Fatal(err)
	}
	closePool(pool)
}
