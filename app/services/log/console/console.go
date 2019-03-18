package console

import (
	"context"
	stdLog "log"
	"os"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/color"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

var stdOut = stdLog.New(os.Stdout, "", 0)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Console"
}

func (s Service) Category() string {
	return "log"
}

func (s Service) Enabled() bool {
	return !env.IsTest()
}

func (s Service) Init() {
	bus.AddEventListener(logDebug)
	bus.AddEventListener(logWarn)
	bus.AddEventListener(logInfo)
	bus.AddEventListener(logError)
}

func logDebug(ctx context.Context, cmd *log.DebugCommand) error {
	writeLog(ctx, log.DEBUG, cmd.Message, cmd.Props)
	return nil
}

func logWarn(ctx context.Context, cmd *log.WarnCommand) error {
	writeLog(ctx, log.WARN, cmd.Message, cmd.Props)
	return nil
}

func logInfo(ctx context.Context, cmd *log.InfoCommand) error {
	writeLog(ctx, log.INFO, cmd.Message, cmd.Props)
	return nil
}

func logError(ctx context.Context, cmd *log.ErrorCommand) error {
	if cmd.Err != nil {
		writeLog(ctx, log.ERROR, cmd.Err.Error(), cmd.Props)
	} else {
		writeLog(ctx, log.ERROR, "nil", cmd.Props)
	}
	return nil
}

func writeLog(ctx context.Context, level log.Level, message string, props log.Props) {
	props = log.GetProps(ctx).Merge(props)
	if log.CurrentLevel > level {
		return
	}

	message = log.Parse(message, props, true)
	contextID := props[log.PropertyKeyContextID]
	if contextID == nil {
		contextID = "???"
	}

	tag := props[log.PropertyKeyTag]
	if tag == nil {
		tag = "???"
	}
	stdOut.Printf("%s [%s] [%s] [%s] %s\n", colorizeLevel(level), time.Now().Format(time.RFC3339), tag, contextID, message)
}

func colorizeLevel(level log.Level) string {
	var colorFunc func(interface{}) string

	switch level {
	case log.DEBUG:
		colorFunc = color.Magenta
	case log.INFO:
		colorFunc = color.Blue
	case log.WARN:
		colorFunc = color.Yellow
	case log.ERROR:
		colorFunc = color.Red
	default:
		colorFunc = color.Red
	}

	return color.Bold(colorFunc(level.String()))
}
