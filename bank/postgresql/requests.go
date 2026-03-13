package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

type Request struct {
	id              uint64
	expectedSuccess bool
	fromID          uint64
	toID            uint64
	amount          uint64
}

func readRequests(inputFile string) (uint64, uint64, []Request, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to open the input file: %w", err)
	}
	var requests []Request
	scanner := bufio.NewScanner(f)
	// The first line is the number of accounts and initial balance
	if !scanner.Scan() {
		return 0, 0, nil, fmt.Errorf("the input file is empty")
	}
	firstLine := strings.Fields(strings.TrimSpace(scanner.Text()))
	if len(firstLine) < 2 {
		return 0, 0, nil, fmt.Errorf("first line: expected two fields, got %d", len(firstLine))
	}
	accounts, err := strconv.ParseUint(firstLine[0], 10, 64)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("first line: bad number of accounts: %w", err)
	}
	initialBalance, err := strconv.ParseUint(firstLine[1], 10, 64)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("first line: bad initial balance: %w", err)
	}
	// Read each request
	for lineNum := 0; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return 0, 0, nil, fmt.Errorf("line %d: expected five fields, got %d", lineNum, len(fields))
		}
		id, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("line %d: bad id: %w", lineNum, err)
		}
		expectedSuccess := fields[1] == "1"
		fromAccount, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("line %d: bad from_account: %w", lineNum, err)
		}
		toAccount, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("line %d: bad to_account: %w", lineNum, err)
		}
		amount, err := strconv.ParseInt(fields[4], 10, 64)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("line %d: bad amount: %w", lineNum, err)
		}
		requests = append(requests, Request{
			id:              id,
			expectedSuccess: expectedSuccess,
			fromID:          uint64(fromAccount),
			toID:            uint64(toAccount),
			amount:          uint64(amount),
		})
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, nil, fmt.Errorf("reading input file: %w", err)
	}
	fmt.Printf("Read %d requests from %s\n", len(requests), inputFile)
	if err := f.Close(); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to close the input file: %w", err)
	}
	return accounts, initialBalance, requests, nil
}

type RequestProducer struct {
	requests []Request
	index    atomic.Uint64
	looping  bool
}

func NewRequestProducer(requests []Request, looping bool) *RequestProducer {
	return &RequestProducer{
		requests: requests,
		looping:  looping,
	}
}

func (rp *RequestProducer) Next() (*Request, error) {
	i := rp.index.Add(1) - 1
	totalRequests := uint64(len(rp.requests))
	if i >= totalRequests && !rp.looping {
		return nil, fmt.Errorf("no more requests")
	}
	return &rp.requests[i%totalRequests], nil
}
