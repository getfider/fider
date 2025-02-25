//go:build windows
// +build windows

package cmd

import (
	"os"

	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
)

var extraSignals = []os.Signal{}

func handleExtraSignal(s os.Signal, e *web.Engine) int {
	return -1
}
