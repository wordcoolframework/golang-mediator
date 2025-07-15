package builders

import (
	"github.com/wordcoolframework/golang-mediator/pkg/mediator"
	"github.com/wordcoolframework/golang-mediator/pkg/mediator/contracts"
)

type Builder struct {
	mediator *mediator.Mediator
}

func NewBuilder() *Builder {
	return &Builder{
		mediator: mediator.New(),
	}
}

func (b *Builder) UseBehavior(behavior mediator.Behavior) *Builder {
	b.mediator.UseBehavior(behavior)
	return b
}

func (b *Builder) Register(handler interface{}) *Builder {
	b.mediator.Register(handler)
	return b
}

func (b *Builder) RegisterEventHandler(event contracts.Event, handler contracts.IEventHandler) *Builder {
	b.mediator.RegisterEventHandler(event, handler)
	return b
}

func (b *Builder) PublishEvent(event contracts.Event) *Builder {
	b.mediator.PublishEvent(event)
	return b
}

func (b *Builder) UseRabbitMQ(url string) *Builder {
	err := b.mediator.UseRabbitMQ(url)
	if err != nil {
		return nil
	}
	return b
}

func (b *Builder) PublishEventToQueue(event contracts.Event) *Builder {
	b.mediator.PublishEvent(event)
	return b
}

func (b *Builder) Provide(dep interface{}) *Builder {
	b.mediator.Provide(dep)
	return b
}

func (b *Builder) Build() *mediator.Mediator {
	return b.mediator
}
