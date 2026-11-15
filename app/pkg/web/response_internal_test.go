package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
)

func TestIsClientGone(t *testing.T) {
	gone := []error{
		context.Canceled,
		fmt.Errorf("write: %w", context.Canceled),
		&net.OpError{Op: "write", Net: "tcp", Err: os.NewSyscallError("write", syscall.EPIPE)},
		&net.OpError{Op: "write", Net: "tcp", Err: os.NewSyscallError("write", syscall.ECONNRESET)},
		errors.New("http2: stream closed"),
		errors.New("client disconnected"),
	}
	for _, err := range gone {
		if !isClientGone(err) {
			t.Errorf("expected client gone for %v", err)
		}
	}

	notGone := []error{
		errors.New("disk on fire"),
		&net.OpError{Op: "write", Net: "tcp", Err: os.NewSyscallError("write", syscall.ENOSPC)},
		context.DeadlineExceeded,
	}
	for _, err := range notGone {
		if isClientGone(err) {
			t.Errorf("expected real error for %v", err)
		}
	}
}
