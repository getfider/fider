// +build !windows

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/getfider/fider/app/pkg/web"
)

var extraSignals = []os.Signal{syscall.SIGUSR1}

func handleExtraSignal(s os.Signal, e *web.Engine) int {
	switch s {
	case syscall.SIGUSR1:
		println("SIGUSR1 received")
		println("Dumping process status")
		buf := new(bytes.Buffer)
		_ = pprof.Lookup("goroutine").WriteTo(buf, 1)
		_ = pprof.Lookup("heap").WriteTo(buf, 1)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("# Worker Queue: %d\n", e.Worker().Length()))
		buf.WriteString(fmt.Sprintf("# Num Goroutines: %d\n", runtime.NumGoroutine()))
		println(buf.String())
	}
	return -1
}
