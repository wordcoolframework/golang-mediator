package builders

import "github.com/wordcoolframework/golang-mediator/pkg/mediator"

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

func (b *Builder) Build() *mediator.Mediator {
	return b.mediator
}
