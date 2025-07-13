package container

import (
	"reflect"
)

type Container struct {
	dependencies map[reflect.Type]interface{}
}

func NewContainer() *Container {
	return &Container{
		dependencies: make(map[reflect.Type]interface{}),
	}
}

func (c *Container) Provide(instance interface{}) {
	t := reflect.TypeOf(instance)
	c.dependencies[t] = instance
}

func (c *Container) Resolve(t reflect.Type) (interface{}, bool) {
	v, ok := c.dependencies[t]
	return v, ok
}
