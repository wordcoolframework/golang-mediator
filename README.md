![Go Mediator Banner](pkg/assets/golang-mediator.png)

# Golang Mediator - CQRS Pattern Implementation

Overview

- ✅ Auto-discovery of handlers based on naming convention

- ✅ Behavior pipeline (similar to middleware)

- ✅ Command/Query separation (CQRS pattern)

- ✅ Dependency Injection support

- ✅ Event handling system

- ✅ RabbitMQ integration for event publishing

- ✅ Clean architecture

- ✅ Simple and testable design

- ✅ Chainable builder pattern for configuration

---

### simple Project Structure
```
/app
  /Commands
  /Queries
  /CommandHandlers
  /QueryHandlers
  /Events
  /EventHandlers
  /Services
```

## 📦 Installation

```bash
go get github.com/wordcoolframework/golang-mediator@latest
```

## 🚀 Quick Start

### 1. Define a Query


```go
// app/Queries/GetUserQuery.go
package Queries

type GetUserQuery struct {
    UserID uint
}
```

### 2. Create the Handler

```go
// app/QueryHandlers/GetUserQueryHandler.go
package QueryHandlers

type GetUserQueryHandler struct{}

func (h *GetUserQueryHandler) Handle(query contracts.Query) (interface{}, error) {
    q := query.(Queries.GetUserQuery)
    return map[string]interface{}{
        "id":   q.UserID,
        "name": "Arash 1380",
    }, nil
}
```

### 3. Use Mediator in Main

```go
m := mediator.New()

m.Register(QueryHandlers.GetUserQueryHandler{})

app.Get("/user/:id", func(c *fiber.Ctx) error {

    id, _ := strconv.Atoi(c.Params("id"))
    query := Queries.GetUserQuery{ID: id}
    
    res, err := m.Send(query)
    
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }
    
    return c.JSON(res)
})
```

### 4. Behaviors (Middleware-like Plugins)

```go
m.UseBehavior(behaviors.LogBehavior)
m.UseBehavior(behaviors.TimerBehavior)
```


### 5. Builder (Use Builder Method Like Chainable)

```go
m := builders.NewBuilder().
    UseBehavior(behaviors.LogBehavior).
    Register(&QueryHandlers.GetUserQueryHandler{}).
    Build()

app.Get("/user/:id", func(c *fiber.Ctx) error {
    
    id, _ := strconv.Atoi(c.Params("id"))
    query := Queries.GetUserQuery{ID: uint(id)}
    
    res, err := m.Send(query)
    
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }
    
        return c.JSON(res)
    })

app.Listen(":3000")
```

### 6. DI (Inject Services)
```go

package Services

type UserService struct{}

func (u *UserService) GetUser(username string) string {
    return "Hello " + username
}


package QueryHandlers

type GetUserQueryHandler struct {
    userService Services.UserService
}

func (h *GetUserQueryHandler) Handle(q Queries2.GetUserQuery) (map[string]string, error) {
    return map[string]string{
    "id":                fmt.Sprintf("%d", q.ID),
    "name":              "Arash Narimani",
    "user-service-data": h.userService.GetUser("arash"),
    }, nil
}


package main

// use Provide(&Services.UserService{}).
m := builders.NewBuilder().
    UseBehavior(behaviors.LogBehavior).
    Register(&QueryHandlers.GetUserQueryHandler{}).
    Provide(&Services.UserService{}).
    Build()
```

### 7. Event (Handle Events With EventHandlers)

```go
package Events

type UserCreatedEvent struct {
	UserID uint
	Email  string
}

func (e UserCreatedEvent) EventName() string {
	return "UserCreatedEvent"
}


package EventHandlers

type UserCreatedHandler struct{}

func (h *UserCreatedHandler) Handle(e contracts.Event) error {
	ev := e.(Events.UserCreatedEvent)
	fmt.Println("Send to Broker:", ev.UserID, ev.Email)
	return nil
}

package main

// use RegisterEventHandler(Events.UserCreatedEvent{}, &EventHandlers.UserCreatedHandler{}).

m := builders.NewBuilder().
    UseBehavior(behaviors.LogBehavior).
    Register(&QueryHandlers.GetUserQueryHandler{}).
    RegisterEventHandler(Events.UserCreatedEvent{}, &EventHandlers.UserCreatedHandler{}).
    Provide(&Services.UserService{}).
    Build()

app.Get("/user/:id", func(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))
	query := Queries.GetUserQuery{ID: uint(id)}

	res, err := m.Send(query)

	ev := Events.UserCreatedEvent{UserID: 10, Email: "arash@gmail.com"}
	m.PublishEvent(ev)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(res)
})

```

### 7. RabbitMq (Publish Data To Message Broker)

```go
package main

// use | UseRabbitMQ("amqp://guest:guest@localhost:5672/").
// publish data | errPublishToQueue := m.PublishEventToQueue(useEvent)

m := builders.NewBuilder().
    UseBehavior(behaviors.LogBehavior).
    Register(&QueryHandlers.GetUserQueryHandler{}).
    RegisterEventHandler(Events.UserCreatedEvent{}, &EventHandlers.UserCreatedHandler{}).
    Provide(&Services.UserService{}).
    UseRabbitMQ("amqp://guest:guest@localhost:5672/").
    Build()

app.Get("/user/:id", func(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))
	query := Queries.GetUserQuery{ID: uint(id)}

	res, err := m.Send(query)

	useEvent := Events.UserCreatedEvent{UserID: 10, Email: "arash@gmail.com"}
	errPublishToQueue := m.PublishEventToQueue(useEvent)

	if errPublishToQueue != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(res)
})

app.Listen(":3000")
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.