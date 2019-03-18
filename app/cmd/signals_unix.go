// +build !windows

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

var extraSignals = []os.Signal{syscall.SIGUSR1}

func handleExtraSignal(s os.Signal, e *web.Engine, settings *models.SystemSettings) int {
	switch s {
	case syscall.SIGUSR1:
		println("SIGUSR1 received")
		println("Dumping process status")
		buf := new(bytes.Buffer)
		_ = pprof.Lookup("goroutine").WriteTo(buf, 1)
		_ = pprof.Lookup("heap").WriteTo(buf, 1)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("# FIDER v%s\n", settings.Version))
		buf.WriteString(fmt.Sprintf("# BuildTime: %s\n", settings.BuildTime))
		buf.WriteString(fmt.Sprintf("# Compiler: %s\n", settings.Compiler))
		buf.WriteString(fmt.Sprintf("# Environment: %s\n", settings.Environment))
		buf.WriteString(fmt.Sprintf("# Worker Queue: %d\n", e.Worker().Length()))
		buf.WriteString(fmt.Sprintf("# Num Goroutines: %d\n", runtime.NumGoroutine()))
		println(buf.String())
	}
	return -1
}
