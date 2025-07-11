package behaviors

import (
	"fmt"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
	"time"
)

func TimerBehavior(req contracts.Request, next func(request contracts.Request) (any, error)) (any, error) {
	start := time.Now()
	res, err := next(req)
	fmt.Printf("%T took %s\n", req, time.Since(start))
	return res, err
}
