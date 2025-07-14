package events

import (
	"fmt"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
)

type EventBus struct {
	registry *eventHandlerRegistry
}

func NewEventBus() *EventBus {
	return &EventBus{
		registry: newEventHandlerRegistry(),
	}
}

func (eb *EventBus) RegisterEventHandler(event contracts.Event, handler contracts.IEventHandler) {
	eb.registry.register(event, handler)
}

func (eb *EventBus) Publish(event contracts.Event) {
	handlers := eb.registry.getHandlers(event)
	for _, handler := range handlers {
		go func(h contracts.IEventHandler) {
			err := h.Handle(event)
			if err != nil {
				fmt.Println("Error handling event:", err)
			}
		}(handler)
	}
}
