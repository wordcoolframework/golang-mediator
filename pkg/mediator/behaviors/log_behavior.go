package behaviors

import (
	"fmt"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
)

func LogBehavior(req contracts.Request, next func(request contracts.Request) (any, error)) (any, error) {
	fmt.Printf("Handling %T\n", req)
	res, err := next(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result : %#v\n", res)
	}
	return res, err
}
