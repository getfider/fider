package sql

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "SQL"
}

func (s Service) Category() string {
	return "log"
}

func (s Service) Enabled() bool {
	return !env.IsTest() && env.Config.Log.Sql
}

func (s Service) Init() {
	bus.AddListener(logDebug)
	bus.AddListener(logWarn)
	bus.AddListener(logInfo)
	bus.AddListener(logError)
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

	message = log.Parse(message, props, false)
	tag := props[log.PropertyKeyTag]
	if tag == nil {
		tag = "???"
	}
	delete(props, log.PropertyKeyTag)

	go func() {
		_, _ = dbx.Connection().Exec(
			"INSERT INTO logs (tag, level, text, created_at, properties) VALUES ($1, $2, $3, $4, $5)",
			tag, level.String(), message, time.Now(), props,
		)
	}()
}
