package file

import (
	"context"
	"encoding/json"
	stdLog "log"
	"os"
	"path/filepath"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

var stdOut *stdLog.Logger

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "File"
}

func (s Service) Category() string {
	return "log"
}

func (s Service) Enabled() bool {
	return !env.IsTest() && env.Config.Log.File
}

func (s Service) Init() {
	bus.AddListener(logDebug)
	bus.AddListener(logWarn)
	bus.AddListener(logInfo)
	bus.AddListener(logError)

	// Get configured output file and directory
	file := env.Config.Log.OutputFile
	dir := filepath.Dir(file)

	// Ensure directory exists
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			panic(err)
		}
	}

	// Open (or create) output file
	output, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	stdOut = stdLog.New(output, "", 0)
}

func logDebug(ctx context.Context, c *cmd.LogDebug) {
	writeLog(ctx, log.DEBUG, c.Message, c.Props)
}

func logWarn(ctx context.Context, c *cmd.LogWarn) {
	writeLog(ctx, log.WARN, c.Message, c.Props)
}

func logInfo(ctx context.Context, c *cmd.LogInfo) {
	writeLog(ctx, log.INFO, c.Message, c.Props)
}

func logError(ctx context.Context, c *cmd.LogError) {
	if c.Err != nil {
		writeLog(ctx, log.ERROR, c.Err.Error(), c.Props)
	} else if c.Message != "" {
		writeLog(ctx, log.ERROR, c.Message, c.Props)
	} else {
		writeLog(ctx, log.ERROR, "nil", c.Props)
	}
}

func writeLog(ctx context.Context, level log.Level, message string, props dto.Props) {
	if !log.IsEnabled(level) {
		return
	}

	props = log.GetProperties(ctx).Merge(props)
	props["Level"] = level.String()
	props["Message"] = log.Parse(message, props, false)
	props["Timestamp"] = time.Now().Format(time.RFC3339)
	if props[log.PropertyKeyTag] == nil {
		props[log.PropertyKeyTag] = "???"
	}

	if env.Config.Log.Structured {
		_ = json.NewEncoder(stdOut.Writer()).Encode(props)
		return
	}

	stdOut.Printf("%s [%s] [%s] %s\n", level, props["Timestamp"], props[log.PropertyKeyTag], props["Message"])
}
