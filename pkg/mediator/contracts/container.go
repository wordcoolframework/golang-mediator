package contracts

import "reflect"

type IContainer interface {
	Provide(instance interface{})
	Resolve(t reflect.Type) (interface{}, bool)
}
