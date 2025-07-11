package mediator

import "github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"

type Behavior func(req contracts.Request, next func(request contracts.Request) (any, error)) (any, error)

func wrapBehavior(b Behavior, next func(request contracts.Request) (any, error)) func(request contracts.Request) (any, error) {
	return func(r contracts.Request) (any, error) {
		return b(r, next)
	}
}
