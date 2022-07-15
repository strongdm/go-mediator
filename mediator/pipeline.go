package mediator

import "context"

type (
	Behavior  func(context.Context, Message, Next) (interface{}, error)
	Behaviors []Behavior
)

func (b Behaviors) reverseApply(fn func(Behavior)) {
	for i := len(b) - 1; i >= 0; i-- {
		fn(b[i])
	}
}

type Next func(ctx context.Context) (interface{}, error)

type Pipeline func(context.Context, Message) (interface{}, error)

type Option func(pCtx *PipelineContext) error

func (p Pipeline) empty() bool { return p == nil }

type PipelineContext struct {
	behaviors Behaviors
	pipeline  Pipeline
	handlers  map[string]func() RequestHandler
}

func newPipelineContext(opts ...Option) (*PipelineContext, error) {
	ctx := PipelineContext{
		handlers: make(map[string]func() RequestHandler),
	}
	for _, opt := range opts {
		if err := opt(&ctx); err != nil {
			return nil, err
		}
	}
	return &ctx, nil
}

func WithBehaviour(behavior PipelineBehaviour) Option {
	return func(pCtx *PipelineContext) error {
		return pCtx.useBehavior(behavior)
	}
}

func WithBehaviourFunc(fn func(context.Context, Message, Next) (interface{}, error)) Option {
	return func(pCtx *PipelineContext) error {
		return pCtx.use(fn)
	}
}

func WithHandler(req Message, rh RequestHandler) Option {
	return func(pCtx *PipelineContext) error {
		f := func() RequestHandler { return rh }
		return pCtx.registerHandlerFunc(req, f)
	}
}

func WithHandlerFunc(req Message, rhf func() RequestHandler) Option {
	return func(pCtx *PipelineContext) error {
		return pCtx.registerHandlerFunc(req, rhf)
	}
}

func (p *PipelineContext) useBehavior(behavior PipelineBehaviour) error {
	if behavior == nil {
		return ErrInvalidArg
	}
	return p.use(behavior.Process)
}

func (p *PipelineContext) use(call func(context.Context, Message, Next) (interface{}, error)) error {
	if call == nil {
		return ErrInvalidArg
	}
	p.behaviors = append(p.behaviors, call)
	return nil
}

func (p *PipelineContext) registerHandlerFunc(req Message, h func() RequestHandler) error {
	if req == nil || h == nil {
		return ErrInvalidArg
	}
	key := req.Key()
	p.handlers[key] = h

	return nil
}
