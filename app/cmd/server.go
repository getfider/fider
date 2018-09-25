package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
)

//RunServer starts the Fider Server
//Returns an exitcode, 0 for OK and 1 for ERROR
func RunServer(settings *models.SystemSettings) int {
	exit := RunMigrate()
	if exit != 0 {
		return exit
	}

	e := routes(web.New(settings))

	go e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))
	return listenSignals(e, settings)
}

func listenSignals(e *web.Engine, settings *models.SystemSettings) int {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, append([]os.Signal{syscall.SIGTERM, syscall.SIGINT}, extraSignals...)...)
	for {
		s := <-signals
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			err := e.Stop()
			if err != nil {
				e.Logger().Error(errors.Wrap(err, "failed to stop fider"))
				return 1
			}
			return 0
		default:
			ret := handleExtraSignal(s, e, settings)
			if ret >= 0 {
				return ret
			}
		}
	}
}
