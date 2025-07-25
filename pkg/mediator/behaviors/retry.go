package behaviors

import (
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
	"time"
)

var MaxRetries int = 3
var Delay time.Duration = 500 * time.Millisecond

func Retry(req contracts.Request, next func(request contracts.Request) (any, error)) (any, error) {
	var result any
	var err error

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		result, err = next(req)
		if err == nil {
			return result, nil
		}
		if attempt < MaxRetries {
			time.Sleep(Delay)
		}
	}
	return result, err
}
