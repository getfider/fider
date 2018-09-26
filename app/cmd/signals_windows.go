// +build windows

package cmd

import (
	"os"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

var extraSignals = []os.Signal{}

func handleExtraSignal(s os.Signal, e *web.Engine, settings *models.SystemSettings) int {
	return -1
}
