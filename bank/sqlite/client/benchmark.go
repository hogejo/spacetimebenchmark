package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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

func isAcceptedError(response string) bool {
	return strings.Contains(response, "same_account") ||
		strings.Contains(response, "insufficient_funds") ||
		strings.Contains(response, "account_not_found")
}

func worker(wg *sync.WaitGroup, counters *Counters, needToStop *atomic.Bool, conn *Connection, producer *RequestProducer, retries int) {
	defer wg.Done()
	for !needToStop.Load() {
		request, err := producer.Next()
		if err != nil {
			return
		}
		var response string
		acceptedError := false
		for attempt := 1; attempt <= retries; attempt++ {
			response, err = conn.sendTransfer(request.fromID, request.toID, request.amount)
			if err != nil {
				// Connection error
				if attempt < retries {
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				}
				continue
			}
			if response == "OK" && request.expectedSuccess {
				break
			}
			acceptedError = isAcceptedError(response)
			if response != "OK" && !request.expectedSuccess && acceptedError {
				break
			}
			if attempt < retries {
				time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			}
		}
		if err != nil {
			fmt.Printf("connection error: from=%d to=%d amount=%d: %v\n", request.fromID, request.toID, request.amount, err)
			counters.errors.Add(1)
			continue
		}
		if response == "OK" && request.expectedSuccess {
			counters.success.Add(1)
			continue
		}
		if response != "OK" && !request.expectedSuccess && acceptedError {
			counters.failures.Add(1)
			continue
		}
		fmt.Printf("unexpected failure: from=%d to=%d amount=%d: %s\n", request.fromID, request.toID, request.amount, response)
		counters.errors.Add(1)
	}
}

func runBenchmark(config Config, requests []Request) {
	connections := config.maximumConnections
	counters := Counters{}
	needToStop := atomic.Bool{}
	needToStop.Store(false)
	producer := NewRequestProducer(requests, true)
	var wg sync.WaitGroup
	wg.Add(connections)
	for i := 0; i < connections; i++ {
		conn, err := newConnection(config.serverAddr)
		if err != nil {
			fmt.Printf("failed to connect worker %d: %v\n", i, err)
			wg.Done()
			continue
		}
		go func(c *Connection) {
			defer c.Close()
			worker(&wg, &counters, &needToStop, c, producer, config.retries)
		}(conn)
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
