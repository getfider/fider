package log

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
)

var CurrentLevel = parseLevel(strings.ToUpper(env.Config.Log.Level))

func IsEnabled(level Level) bool {
	return CurrentLevel <= level
}

func Debug(ctx context.Context, message string) {
	bus.Publish(ctx, &cmd.LogDebug{Message: message})
}

func Debugf(ctx context.Context, message string, props dto.Props) {
	bus.Publish(ctx, &cmd.LogDebug{Message: message, Props: props})
}

func Info(ctx context.Context, message string) {
	bus.Publish(ctx, &cmd.LogInfo{Message: message})
}

func Infof(ctx context.Context, message string, props dto.Props) {
	bus.Publish(ctx, &cmd.LogInfo{Message: message, Props: props})
}

func Warn(ctx context.Context, message string) {
	bus.Publish(ctx, &cmd.LogWarn{Message: message})
}

func Warnf(ctx context.Context, message string, props dto.Props) {
	bus.Publish(ctx, &cmd.LogWarn{Message: message, Props: props})
}

func Error(ctx context.Context, err error) {
	bus.Publish(ctx, &cmd.LogError{Err: err})
}

func Errorf(ctx context.Context, message string, props dto.Props) {
	bus.Publish(ctx, &cmd.LogError{Message: message, Props: props})
}
