package bus_test

import (
	"context"

	"github.com/getfider/fider/app/pkg/bus"
)

type SayHelloCommand struct {
	Name   string
	Result string
}

type key int

const (
	GreetingKey key = iota
)

type GreeterService struct {
}

func (s GreeterService) Name() string {
	return "Greeter"
}

func (s GreeterService) Category() string {
	return "greeter"
}

func (s GreeterService) Enabled() bool {
	return true
}

func (s GreeterService) Init() {
	bus.AddHandler(Greet)
}

var getGreeting = func(ctx context.Context) string {
	return ctx.Value(GreetingKey).(string)
}

func Greet(ctx context.Context, cmd *SayHelloCommand) error {
	cmd.Result = getGreeting(ctx) + " " + cmd.Name
	return nil
}

type BetterGreeterService struct {
}

func (s BetterGreeterService) Name() string {
	return "BetterGreeter"
}

func (s BetterGreeterService) Category() string {
	return "greeter"
}

func (s BetterGreeterService) Enabled() bool {
	return true
}

func (s BetterGreeterService) Init() {
	bus.AddHandler(SayHello)
}

func SayHello(ctx context.Context, cmd *SayHelloCommand) error {
	cmd.Result = "Hello " + cmd.Name
	return nil
}
