package main

import (
	"fmt"
	"log"
	"math/rand/v2"
)

func main() {
	log.SetFlags(0)
	config := parseConfig()
	printSummary(&config)
	// Pre-decide which requests are failures
	isFailingRequest := buildFailureSchedule(config.totalRequests, config.failingRequests)
	// Generate requests
	rw := NewRequestsWriter(config.output)
	rw.WriteHeader(config)
	rng := rand.New(rand.NewPCG(12345, 67890))
	getAccount := func() uint64 { return rng.Uint64N(uint64(float64(config.accounts) * config.coverage)) }
	if config.zipf {
		zipf := rand.NewZipf(rng, config.coverage, 1.0, uint64(config.accounts-1))
		getAccount = func() uint64 { return zipf.Uint64() }
	}
	lastPercent := -1
	for i := uint64(0); i < config.totalRequests; i++ {
		percent := int((i + 1) * 100 / config.totalRequests)
		if percent != lastPercent {
			fmt.Printf("Generating requests ... %d%%\r", percent)
			lastPercent = percent
		}
		fromID := getAccount()
		toID := getAccount()
		for fromID == toID {
			toID = getAccount()
		}
		amount := uint64(rand.IntN(100) + 1)
		isSuccessful := true
		if isFailingRequest[i] {
			isSuccessful = false
			fromID, toID, amount = applyFailure(fromID, toID, amount, config.accounts)
		}
		rw.WriteLine(i, isSuccessful, fromID, toID, amount)
	}
	fmt.Println()
	rw.Close()
}
