package behaviors

import (
	"errors"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
	"sync"
	"time"
)

var (
	ErrCircuitOpen   = errors.New("circuit breaker is open")
	FailureThreshold = 5
	ResetTimeout     = 30 * time.Second
	failCount        = 0
	state            = "closed" // "closed", "open", "half-open"
	lastFail         time.Time
	mutex            sync.Mutex
)

func CircuitBreaker(req contracts.Request, next func(contracts.Request) (any, error)) (any, error) {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()

	if state == "open" {
		if now.Sub(lastFail) >= ResetTimeout {
			state = "half-open"
		} else {
			return nil, ErrCircuitOpen
		}
	}

	result, err := next(req)

	if err != nil {
		failCount++
		lastFail = now
		if failCount >= FailureThreshold {
			state = "open"
		}
		return result, err
	}

	failCount = 0
	state = "closed"
	return result, nil
}
