package mediator_test

import (
	"context"
	"testing"

	"github.com/TheFellow/go-mediator/mediator"
	"github.com/stretchr/testify/assert"
)

func TestMediator_should_dispatch_msg_when_send(t *testing.T) {
	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	_, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, cmd, handler.captured)
}

func TestMediator_should_return_handler_result(t *testing.T) {
	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	result, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, cmd.name, result)
}

func TestMediator_should_execute_behavior_when_send(t *testing.T) {
	var got mediator.Message
	behavior := func(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
		got = msg
		return next(ctx)
	}

	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithBehaviourFunc(behavior),
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	_, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, cmd, got)
}

func TestMediator_with_behavior_func_should_return_handler_result(t *testing.T) {
	passThru := func(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
		return next(ctx)
	}

	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithBehaviourFunc(passThru),
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	result, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, cmd.name, result)
}

func TestMediator_with_behavior_should_return_handler_result(t *testing.T) {

	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithBehaviour(PassThruBehavior{}),
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	result, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, cmd.name, result)
}

func TestMediator_with_behavior_func_can_alter_handler_result(t *testing.T) {
	theResultIs42 := func(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
		_, err := next(ctx)
		return 42, err
	}

	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithBehaviourFunc(theResultIs42),
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	result, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestMediator_with_behavior_can_alter_handler_result(t *testing.T) {

	cmd := &fakeCommand{
		name: "Amsterdam",
	}
	handler := &fakeCommandHandler{}

	m, _ := mediator.New(
		mediator.WithBehaviour(FortyTwoBehavior{}),
		mediator.WithHandler(&fakeCommand{}, handler),
	)

	result, err := m.Send(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

type fakeCommand struct {
	name string
}

func (*fakeCommand) Key() string { return "fakeCommand" }

type fakeCommandHandler struct {
	captured mediator.Message
}

func (f *fakeCommandHandler) Handle(_ context.Context, msg mediator.Message) (interface{}, error) {
	f.captured = msg
	cmd := msg.(*fakeCommand)
	return cmd.name, nil
}

type PassThruBehavior struct{}

func (p PassThruBehavior) Process(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
	return next(ctx)
}

type FortyTwoBehavior struct{}

func (p FortyTwoBehavior) Process(ctx context.Context, msg mediator.Message, next mediator.Next) (interface{}, error) {
	_, err := next(ctx)
	return 42, err
}
