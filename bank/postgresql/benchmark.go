package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Counters struct {
	success  atomic.Uint64
	failures atomic.Uint64
	errors   atomic.Uint64
}

func (c *Counters) Reset() {
	c.success.Store(0)
	c.failures.Store(0)
	c.errors.Store(0)
}

func isAcceptedError(err error, request *Request) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "same_account") ||
		strings.Contains(msg, "insufficient_funds") ||
		strings.Contains(msg, "account_not_found")
}

func worker(wg *sync.WaitGroup, counters *Counters, needToStop *atomic.Bool, pool *pgxpool.Pool, producer *RequestProducer, retries int) {
	defer wg.Done()
	for !needToStop.Load() {
		request, err := producer.Next()
		if err != nil {
			return
		}
		acceptedError := false
		for attempt := 1; attempt <= retries; attempt++ {
			_, err = pool.Exec(context.Background(), "CALL transfer($1, $2, $3)", request.fromID, request.toID, request.amount)
			if err == nil && request.expectedSuccess {
				break
			}
			acceptedError = isAcceptedError(err, request)
			if err != nil && !request.expectedSuccess && acceptedError {
				break
			}
			if attempt < retries {
				time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			}
		}
		if err == nil && request.expectedSuccess {
			counters.success.Add(1)
			continue
		}
		if err != nil && !request.expectedSuccess && acceptedError {
			counters.failures.Add(1)
			continue
		}
		fmt.Printf("unexpected failure: from=%d to=%d amount=%d: %v\n", request.fromID, request.toID, request.amount, err)
		counters.errors.Add(1)
	}
}

func runBenchmark(pool *pgxpool.Pool, config Config, requests []Request) {
	//connections, err := acquireConnections(pool, config.maximumConnections)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer releaseConnections(connections)
	workers := config.workers
	counters := Counters{}
	needToStop := atomic.Bool{}
	needToStop.Store(false)
	producer := NewRequestProducer(requests, true)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(&wg, &counters, &needToStop, pool, producer, config.retries)
	}
	fmt.Printf("Workers started. Running warmup for %s ...\n", config.warmupDuration)
	time.Sleep(config.warmupDuration)
	counters.Reset()
	fmt.Printf("Running benchmark for %s ...\n", config.duration)
	time.Sleep(config.duration)
	fmt.Println("Completed with a total of:")
	fmt.Printf("  successful transfers:      %d\n", counters.success.Load())
	fmt.Printf("  expected failed transfers: %d\n", counters.failures.Load())
	fmt.Printf("  errors:                    %d\n", counters.errors.Load())
	rate := float64(counters.success.Load()+counters.failures.Load()) / config.duration.Seconds()
	fmt.Printf("  Throughput:                %.2f TPS\n", rate)
	fmt.Printf("Stopping and waiting for workers ...\n")
	needToStop.Store(true)
	wg.Wait()
}
