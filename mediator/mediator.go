package mediator

import "context"

var _ Sender = &Mediator{}

type Mediator struct {
	context *PipelineContext
}

func New(opts ...Option) (*Mediator, error) {
	pctx, err := newPipelineContext(opts...)
	if err != nil {
		return nil, err
	}
	m := &Mediator{
		context: pctx,
	}

	pctx.behaviors.reverseApply(m.pipe)
	return m, nil
}

func (m *Mediator) Send(ctx context.Context, req Message) (interface{}, error) {
	if m.context.pipeline.empty() {
		return m.send(ctx, req)
	}
	return m.context.pipeline(ctx, req)
}

func (m *Mediator) send(ctx context.Context, req Message) (interface{}, error) {
	key := req.Key()
	handlerFunc, ok := m.context.handlers[key]
	if !ok {
		return nil, ErrHandlerNotFound
	}
	handler, err := handlerFunc()
	if err != nil {
		return nil, err
	}

	return handler.Handle(ctx, req)
}

func (m *Mediator) pipe(call Behavior) {
	if m.context.pipeline.empty() {
		m.context.pipeline = m.send
	}
	seed := m.context.pipeline

	m.context.pipeline = func(ctx context.Context, msg Message) (interface{}, error) {
		return call(ctx, msg, func(context.Context) (interface{}, error) { return seed(ctx, msg) })
	}
}
