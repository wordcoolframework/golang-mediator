package mediator

import (
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/exceptions"
	"reflect"
)

type Mediator struct {
	handlers  map[string]interface{}
	behaviors []Behavior
}

func New() *Mediator {
	return &Mediator{
		handlers:  make(map[string]interface{}),
		behaviors: []Behavior{},
	}
}

func (m *Mediator) UseBehavior(b Behavior) {
	m.behaviors = append(m.behaviors, b)
}

func (m *Mediator) Register(handler interface{}) {
	t := reflect.TypeOf(handler)
	if t.Kind() != reflect.Ptr {
		panic("handler must be a pointer")
	}
	name := t.Elem().Name()
	m.handlers[name] = handler
}

func (m *Mediator) Send(req contracts.Request) (any, error) {
	reqType := reflect.TypeOf(req).Name()
	handlerName := reqType + "Handler"

	handler, ok := m.handlers[handlerName]
	if !ok {
		return nil, exceptions.HandlerNotFoundException
	}

	hVal := reflect.ValueOf(handler)
	method := hVal.MethodByName("Handle")
	if !method.IsValid() {
		return nil, exceptions.HandlerNotFoundException
	}

	call := func(r contracts.Request) (any, error) {
		in := []reflect.Value{reflect.ValueOf(r)}
		out := method.Call(in)
		if len(out) != 2 {
			return nil, exceptions.HandlerNotFoundException
		}
		errVal := out[1].Interface()
		if errVal != nil {
			return out[0].Interface(), errVal.(error)
		}
		return out[0].Interface(), nil
	}

	return m.runBehaviors(req, call)
}

func (m *Mediator) runBehaviors(req contracts.Request, final func(request contracts.Request) (any, error)) (any, error) {
	h := final
	for i := len(m.behaviors) - 1; i >= 0; i-- {
		h = wrapBehavior(m.behaviors[i], h)
	}
	return h(req)
}
