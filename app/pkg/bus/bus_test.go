package bus_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/pkg/bus"

	. "github.com/getfider/fider/app/pkg/assert"
)

var GreetingKey = "GreetingKey"

type GreeterService struct {
}

func (s GreeterService) Category() string {
	return "greeter"
}

func (s GreeterService) Enabled() bool {
	return true
}

func (s GreeterService) Init() {
	bus.AddHandler(s, SayHello)
}

type SayHelloCommand struct {
	Name   string
	Result string
}

var getGreeting = func(ctx context.Context) string {
	return ctx.Value(GreetingKey).(string)
}

func SayHello(ctx context.Context, cmd *SayHelloCommand) error {
	cmd.Result = getGreeting(ctx) + " " + cmd.Name
	return nil
}

func TestBus_SimpleMessage(t *testing.T) {
	RegisterT(t)
	bus.Register(GreeterService{})
	bus.Init()
	ctx := context.WithValue(context.Background(), GreetingKey, "Hello")
	cmd := &SayHelloCommand{Name: "Fider"}
	err := bus.Dispatch(ctx, cmd)
	Expect(err).IsNil()
	Expect(cmd.Result).Equals("Hello Fider")
}

func TestBus_MessageIsNotPointer_ShouldPanic(t *testing.T) {
	RegisterT(t)
	bus.Register(GreeterService{})
	bus.Init()

	defer func() {
		if r := recover(); r != nil {
			panicText := (r.(error)).Error()
			Expect(panicText).Equals("'github.com/getfider/fider/app/pkg/bus_test.SayHelloCommand' is not a pointer")
		}
	}()

	cmd := SayHelloCommand{Name: "Fider"}
	err := bus.Dispatch(context.Background(), cmd)
	Expect(err).IsNil()
}
