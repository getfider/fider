package handlers

import (
	"runtime"

	"github.com/getfider/fider/app/pkg/web"
)

//RuntimeStats returns some useful runtime information
func RuntimeStats() web.HandlerFunc {
	return func(c web.Context) error {
		memStats := &runtime.MemStats{}
		runtime.ReadMemStats(memStats)

		return c.Ok(web.Map{
			"goroutines":  runtime.NumGoroutine(),
			"workerQueue": c.Engine().Worker().Length(),
			"heapInMB":    memStats.HeapAlloc / 1048576,
			"stackInMB":   memStats.StackInuse / 1048576,
		})
	}
}
