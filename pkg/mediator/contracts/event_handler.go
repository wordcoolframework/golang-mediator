package contracts

type IEventHandler interface {
	Handle(Event) error
}
