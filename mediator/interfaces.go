package mediator

import "context"

type (
	Sender interface {
		Send(context.Context, Message) (interface{}, error)
	}
	RequestHandler interface {
		Handle(context.Context, Message) (interface{}, error)
	}
	PipelineBehaviour interface {
		Process(context.Context, Message, Next) (interface{}, error)
	}
	Message interface {
		Key() string
	}
)
