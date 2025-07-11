package contracts

type Request interface{}

type Command interface {
	Request
}

type Query interface {
	Request
}

type Handler[req Request, res any] interface {
	Handle(request req) (res, error)
}
