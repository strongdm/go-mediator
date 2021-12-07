
# go-mediator

Simple mediator implementation in go, offering in-process messaging with behaviours.

## Commands

      
```go
type CreateOrderCommand struct {
    Id string `validate:"required,min=10"`
}

func (CreateOrderCommand) Key() string { reflect.TypeOf(CreateOrderCommand{}).Name() }

type CreateOrderCommandHandler struct {}

func NewCreateOrderCommandHandler() CreateOrderCommandHandler {
    return CreateOrderCommandHandler{}
}

func (CreateOrderCommandHandler) Handle(ctx context.Context, msg mediator.Message) (interface{}, error) {
    cmd := msg.(CreateOrderCommand)
    fmt.Println(cmd.Id) // Or something interesting instead...
    return nil, nil
}
```

## Behaviours

***PipelineBehaviour interface implementation usage***

```go
type Logger struct{}

func NewLogger() *Logger { return &Logger{} }

func (l *Logger) Process(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
    log.Println("Pre Process!")
    result, err := next(ctx)
    log.Println("Post Process...")
    return result, err
}

m, err := mediator.New(mediator.WithBehaviour(behaviour.NewLogger()))
```
  

***Func based usage***

```go
m, err := mediator.New(mediator.WithBehaviourFunc(func(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
    log.Println("Pre Process!")
    next(ctx)
    log.Println("Post Process")
    
    return  nil, nil
}))
```
  

## Usages

```go
m, err := mediator.New(
    mediator.WithBehaviour(behaviour.NewLogger()),
    mediator.WithBehaviour(behaviour.NewValidator()),
    mediator.WithHandler(FakeCommand{}, NewFakeCommandCommandHandler(r)),
)

cmd := FakeCommand{   
    Name: "Amsterdam",  
}

ctx := context.Background()

result, err := m.Send(ctx, cmd)
```
  

***Func based usage***

```go
m, err  := mediator.New(
    mediator.WithBehaviourFunc(func(ctx context.Context, cmd mediator.Message, next mediator.Next) (interface{}, error) {
        log.Println("Pre Process - 1!")
        
        next(ctx)
        
        log.Println("Post Process - 1")
        return  nil, nil
}))
```

## Examples

[Simple](https://github.com/strongdm/go-mediator/tree/master/_examples)

TODO: Add Request/Response example
