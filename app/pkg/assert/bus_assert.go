package assert

import (
	"fmt"
	"reflect"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
)

type BusHandlerAssertions struct {
	msg bus.Msg
}

// ExpectHandler is used to assert bus operations
func ExpectHandler(msg bus.Msg) *BusHandlerAssertions {
	return &BusHandlerAssertions{
		msg,
	}
}

// CalledOnce asserts that a particular handler was called once
func (a *BusHandlerAssertions) CalledOnce() bool {
	return a.CalledTimes(1)
}

// CalledTimes asserts that a particular handler was called X number of times
func (a *BusHandlerAssertions) CalledTimes(n int) bool {
	numOfCalls := bus.GetCallCount(a.msg)
	if reflect.DeepEqual(numOfCalls, n) {
		return true
	}

	err := fmt.Errorf("CalledOnce assertion failed. \n Expected: \n\t\t %s\n Actual: \n\t\t %s", describe(n), describe(numOfCalls))
	currentT.Error(errors.StackN(err, 1))
	return false
}
