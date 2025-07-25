package contracts

type HandlerFunc func(request Request) (interface{}, error)

type Behavior interface {
	Execute(next HandlerFunc) HandlerFunc
}
