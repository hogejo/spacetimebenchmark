package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	defaultCoverage        = 1.5
	defaultAccounts        = 1_000_000
	defaultInitialBalance  = 100_000
	defaultTotalRequests   = 3_000_000
	defaultFailingRequests = defaultTotalRequests / 10
	defaultOutput          = "requests"
)

type Config struct {
	zipf            bool
	coverage        float64
	totalRequests   uint64
	failingRequests uint64
	accounts        uint64
	initialBalance  uint64
	output          string
}

func parseConfig() Config {
	config := Config{}
	flag.BoolVar(&config.zipf, "zipf", false, "Use Zipf distribution")
	flag.Float64Var(&config.coverage, "coverage", defaultCoverage, "Ratio of accounts to generate requests for or alpha when Zipf mode is used")
	flag.Uint64Var(&config.totalRequests, "total-requests", defaultTotalRequests, "Total number of requests to generate")
	flag.Uint64Var(&config.failingRequests, "failing-requests", defaultFailingRequests, "Number of failing requests to generate")
	flag.Uint64Var(&config.accounts, "accounts", defaultAccounts, "Number of accounts")
	flag.Uint64Var(&config.initialBalance, "initial-balance", defaultInitialBalance, "Expected initial balance per account")
	flag.StringVar(&config.output, "output", defaultOutput, "Output file")
	flag.Parse()
	if !config.zipf && (config.coverage <= 0 || config.coverage > 1) {
		log.Fatal("The coverage percentage must be between (0,1]!")
	}
	if config.zipf && config.coverage <= 1.0 {
		log.Fatal("The alpha/exponent must be greater than 1.0 for a valid Zipf distribution!")
	}
	if config.totalRequests == 0 {
		log.Fatal("The number of total requests must be greater than zero!")
	}
	if config.accounts == 0 {
		log.Fatal("The number of accounts must be greater than zero!")
	}
	if config.failingRequests > config.totalRequests {
		log.Fatal("The number of failing requests must not be larger than the total!")
	}
	if config.initialBalance == 0 {
		log.Fatal("The initial balance must be greater than zero!")
	}
	return config
}

func printSummary(config *Config) {
	if config.zipf {
		fmt.Printf("Zipf alpha (exponent): %.3f\n", config.coverage)
	} else {
		fmt.Printf("Coverage is %.3f, requests between the first %d accounts.\n", config.coverage, uint64(float64(config.accounts)*config.coverage))
	}
	fmt.Printf("%d accounts with a %d initial balance\n", config.accounts, config.initialBalance)
	successfulRequests := config.totalRequests - config.failingRequests
	fmt.Printf("Generating %d requests into %s\n", config.totalRequests, config.output)
	if config.failingRequests > 0 {
		fmt.Printf("%d (%d%%) of requests will fail\n",
			config.failingRequests, config.failingRequests*100/config.totalRequests)
		fmt.Printf("%d (%d%%) of requests will succeed\n",
			successfulRequests, successfulRequests*100/config.totalRequests)
	}
}
