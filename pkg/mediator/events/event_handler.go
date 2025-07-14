package events

import "github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"

type eventHandlerRegistry struct {
	handlers map[string][]contracts.IEventHandler
}

func newEventHandlerRegistry() *eventHandlerRegistry {
	return &eventHandlerRegistry{
		handlers: make(map[string][]contracts.IEventHandler),
	}
}

func (r *eventHandlerRegistry) register(e contracts.Event, handler contracts.IEventHandler) {
	eventName := e.EventName()
	r.handlers[eventName] = append(r.handlers[eventName], handler)
}

func (r *eventHandlerRegistry) getHandlers(event contracts.Event) []contracts.IEventHandler {
	return r.handlers[event.EventName()]
}
