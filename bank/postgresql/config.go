package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	defaultDatabaseURL    = "postgres://localhost:5432/benchmark?sslmode=disable"
	defaultWarmupDuration = 5 * time.Second
	defaultDuration       = 10 * time.Second
	defaultConnections    = 10
	defaultMaxInflight    = 128
	defaultWorkers        = 10
	defaultRetries        = 5
	defaultRequestsFile   = "requests"
	defaultVerifyOnly     = false
)

type Config struct {
	databaseURL             string
	input                   string
	maximumConnections      int
	maximumInFlightRequests int
	workers                 int
	retries                 int
	duration                time.Duration
	warmupDuration          time.Duration
	verifyOnly              bool
	// Added from the requests file
	accounts       uint64
	initialBalance uint64
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseConfig() (Config, error) {
	config := Config{}
	flag.StringVar(&config.databaseURL, "database-url", getEnvOrDefault("DATABASE_URL", defaultDatabaseURL), "PostgreSQL connection URL")
	flag.StringVar(&config.input, "input", getEnvOrDefault("INPUT_FILE", defaultRequestsFile), "path to requests file")
	flag.IntVar(&config.maximumConnections, "max-connections", defaultConnections, "number of parallel connections")
	flag.IntVar(&config.maximumInFlightRequests, "max-in-flight", defaultMaxInflight, "global maximum in-flight requests")
	flag.IntVar(&config.workers, "workers", defaultWorkers, "number of workers")
	flag.IntVar(&config.retries, "retries", defaultRetries, "number of retries for failed requests")
	flag.DurationVar(&config.duration, "duration", defaultDuration, "benchmark duration")
	flag.DurationVar(&config.warmupDuration, "warmup-duration", defaultWarmupDuration, "warmup duration")
	flag.BoolVar(&config.verifyOnly, "verify-only", defaultVerifyOnly, "only verify the database")
	flag.Parse()
	if config.input == "" {
		config.input = defaultRequestsFile
	}
	if config.databaseURL == "" {
		config.databaseURL = defaultDatabaseURL
	}
	if config.verifyOnly {
		return config, nil
	}
	if config.maximumConnections <= 0 {
		return config, fmt.Errorf("the number of parallel connections must be greater than zero")
	}
	if config.maximumInFlightRequests <= 0 {
		return config, fmt.Errorf("the maximum in-flight requests per batch must be greater than zero")
	}
	if config.workers < 1 {
		return config, fmt.Errorf("the number of workers must be greater than zero")
	}
	if config.retries < 0 {
		return config, fmt.Errorf("the number of retries must be non-negative")
	}
	if config.duration <= 0 {
		return config, fmt.Errorf("the benchmark duration must be greater than zero")
	}
	if config.warmupDuration < 0 {
		return config, fmt.Errorf("the warmup duration must be non-negative")
	}
	return config, nil
}

func printSummary(config Config) {
	if config.verifyOnly {
		fmt.Println("Only verifying the database...")
		return
	}
	fmt.Printf("Benchmark configuration summary:\n")
	fmt.Printf("  database URL:               %s\n", config.databaseURL)
	fmt.Printf("  input file:                 %s\n", config.input)
	if config.accounts > 0 {
		fmt.Printf("    accounts:                 %d\n", config.accounts)
	}
	if config.initialBalance > 0 {
		fmt.Printf("    initial balance:          %d\n", config.initialBalance)
	}
	fmt.Printf("  maximum connections:        %d\n", config.maximumConnections)
	fmt.Printf("  maximum in-flight requests: %d\n", config.maximumInFlightRequests)
	fmt.Printf("  workers:                    %d\n", config.workers)
	fmt.Printf("  retries:                    %d\n", config.retries)
	fmt.Printf("  benchmark duration:         %s\n", config.duration)
	fmt.Printf("  warmup duration:            %s\n", config.warmupDuration)
}
