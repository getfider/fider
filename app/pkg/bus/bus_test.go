package bus_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/pkg/bus"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestBus_SimpleMessage(t *testing.T) {
	RegisterT(t)
	bus.Register(GreeterService{})
	bus.Init()
	ctx := context.WithValue(context.Background(), GreetingKey, "Good Morning")
	cmd := &SayHelloCommand{Name: "Fider"}
	err := bus.Dispatch(ctx, cmd)
	Expect(err).IsNil()
	Expect(cmd.Result).Equals("Good Morning Fider")
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

func TestBus_OverwriteService(t *testing.T) {
	RegisterT(t)

	bus.Register(GreeterService{})
	bus.Register(BetterGreeterService{})
	bus.Init()
	cmd := &SayHelloCommand{Name: "Fider"}
	err := bus.Dispatch(context.Background(), cmd)
	Expect(err).IsNil()
	Expect(cmd.Result).Equals("Hello Fider")
}
