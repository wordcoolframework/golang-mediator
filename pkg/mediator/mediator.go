package mediator

import (
	"errors"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/container"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/events"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/exceptions"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/rabbitmq"
	"reflect"
)

type BehaviorFunc func(req contracts.Request, next func(req contracts.Request) (any, error)) (any, error)

type Mediator struct {
	handlers       map[string]interface{}
	behaviors      []Behavior
	container      *container.Container
	eventBus       *events.EventBus
	rabbitProducer *rabbitmq.Producer
	pipeline       []BehaviorFunc
}

func New() *Mediator {
	return &Mediator{
		handlers:  make(map[string]interface{}),
		behaviors: []Behavior{},
		container: container.NewContainer(),
		eventBus:  events.NewEventBus(),
	}
}

func (m *Mediator) RegisterEventHandler(event contracts.Event, handler contracts.IEventHandler) {
	m.eventBus.RegisterEventHandler(event, handler)
}

func (m *Mediator) PublishEvent(event contracts.Event) {
	m.eventBus.Publish(event)
}

func (m *Mediator) UseBehavior(b BehaviorFunc) {
	m.pipeline = append(m.pipeline, b)
}

func (m *Mediator) UseRabbitMQ(url string) error {
	producer, err := rabbitmq.NewProducer(rabbitmq.Config{URL: url})
	if err != nil {
		return err
	}
	m.rabbitProducer = producer
	return nil
}

func (m *Mediator) Register(handler interface{}) {
	t := reflect.TypeOf(handler)
	if t.Kind() != reflect.Ptr {
		panic("handler must be a pointer")
	}

	handlerValue := reflect.ValueOf(handler).Elem()
	for i := 0; i < handlerValue.NumField(); i++ {
		field := handlerValue.Type().Field(i)
		if dep, ok := m.container.Resolve(field.Type); ok {
			handlerValue.Field(i).Set(reflect.ValueOf(dep))
		}
	}

	name := t.Elem().Name()
	m.handlers[name] = handler
}

func (m *Mediator) Provide(dep interface{}) {
	m.container.Provide(dep)
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

func (m *Mediator) PublishEventToQueue(event contracts.Event) error {
	if m.rabbitProducer == nil {
		return errors.New("rabbitmq not configured")
	}
	queueName := event.EventName()
	return m.rabbitProducer.Publish(queueName, event)
}

func (m *Mediator) runBehaviors(req contracts.Request, final func(request contracts.Request) (any, error)) (any, error) {
	h := final
	for i := len(m.behaviors) - 1; i >= 0; i-- {
		h = wrapBehavior(m.behaviors[i], h)
	}
	return h(req)
}
