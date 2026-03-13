package main

import (
	"log"
	"math"
	"math/rand/v2"
)

type FailureKind int

const (
	SameAccount FailureKind = iota
	FromNonExistent
	ToNonExistent
	InsufficientFunds
)
const numFailureKinds = 4

var nextFailureKind = 0

func buildFailureSchedule(totalRequests, failingRequests uint64) []bool {
	isFailingRequest := make([]bool, totalRequests)
	if totalRequests > math.MaxInt {
		log.Fatal("More than int.MaxValue number of requests are not supported :(")
	}
	for _, idx := range rand.Perm(int(totalRequests))[:failingRequests] {
		isFailingRequest[idx] = true
	}
	return isFailingRequest
}

func applyFailure(fromID, toID, amount, accounts uint64) (uint64, uint64, uint64) {
	failureKind := FailureKind(nextFailureKind)
	switch failureKind {
	case SameAccount:
		toID = fromID
	case FromNonExistent:
		fromID = accounts
	case ToNonExistent:
		toID = accounts
	case InsufficientFunds:
		amount = math.MaxInt
	}
	nextFailureKind = (nextFailureKind + 1) % numFailureKinds
	return fromID, toID, amount
}
