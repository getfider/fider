package sql

import (
	"context"
	"time"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

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
	if log.CurrentLevel > level {
		return
	}

	usingDatabase(ctx, func(db *dbx.Database) {
		props = log.GetProps(ctx).Merge(props)

		message = log.Parse(message, props, false)
		tag := props[log.PropertyKeyTag]
		if tag == nil {
			tag = "???"
		}
		delete(props, log.PropertyKeyTag)

		_, _ = db.Connection().ExecContext(ctx,
			"INSERT INTO logs (tag, level, text, created_at, properties) VALUES ($1, $2, $3, $4, $5)",
			tag, level.String(), message, time.Now(), props,
		)
	})
}

func usingDatabase(ctx context.Context, handler func(db *dbx.Database)) {
	db, ok := ctx.Value(app.DatabaseCtxKey).(*dbx.Database)
	if ok {
		handler(db)
	}
}
