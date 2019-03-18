package log

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
)

var CurrentLevel = parseLevel(strings.ToUpper(env.Config.Log.Level))

func Debug(ctx context.Context, message string) {
	bus.Publish(ctx, &DebugCommand{Message: message})
}

func Debugf(ctx context.Context, message string, props Props) {
	bus.Publish(ctx, &DebugCommand{Message: message, Props: props})
}

func Info(ctx context.Context, message string) {
	bus.Publish(ctx, &DebugCommand{Message: message})
}

func Infof(ctx context.Context, message string, props Props) {
	bus.Publish(ctx, &DebugCommand{Message: message, Props: props})
}

func Warn(ctx context.Context, message string) {
	bus.Publish(ctx, &WarnCommand{Message: message})
}

func Warnf(ctx context.Context, message string, props Props) {
	bus.Publish(ctx, &WarnCommand{Message: message, Props: props})
}

func Error(ctx context.Context, err error) {
	bus.Publish(ctx, &ErrorCommand{Err: err})
}

func Errorf(ctx context.Context, message string, props Props) {
	bus.Publish(ctx, &ErrorCommand{Message: message, Props: props})
}

type DebugCommand struct {
	Message string
	Props   Props
}

type ErrorCommand struct {
	Err     error
	Message string
	Props   Props
}

type WarnCommand struct {
	Message string
	Props   Props
}

type InfoCommand struct {
	Message string
	Props   Props
}
