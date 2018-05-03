package assert

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/getfider/fider/app/pkg/errors"
)

func mustBeFunction(v interface{}) {
	if reflect.TypeOf(v).Kind() != reflect.Func {
		panic("Value is not a function")
	}
}

//Assertions is the entry point for the assersions
type Assertions struct {
	t *testing.T
}

//New creates a new Assertions entry point for further reuse
func New(t *testing.T) *Assertions {
	return &Assertions{t}
}

//Fail the current test case with given message
func Fail(msg string, args ...interface{}) {
	if currentT == nil {
		panic("Did you forget to call RegisterT(t)?")
	}
	currentT.Errorf(msg, args)
}

var currentT *testing.T

func RegisterT(t *testing.T) {
	currentT = t
}

type AnyAssertions struct {
	actual interface{}
}

func Expect(actual interface{}) *AnyAssertions {
	if currentT == nil {
		panic("Did you forget to call RegisterT(t)?")
	}
	return &AnyAssertions{
		actual: actual,
	}
}

func (a *AnyAssertions) Equals(expected interface{}) bool {
	if reflect.DeepEqual(expected, a.actual) {
		return true
	}
	err := fmt.Errorf("Equals assertion failed. \n Expected: \n\t\t [%s] %v \n Actual: \n\t\t [%s] %v", reflect.ValueOf(expected).Type(), expected, reflect.ValueOf(a.actual).Type(), a.actual)
	currentT.Error(errors.StackN(err, 1))
	return false
}

func (a *AnyAssertions) ContainsSubstring(substr string) bool {
	if strings.Contains(a.actual.(string), substr) {
		return true
	}
	err := fmt.Errorf("ContainsSubstring assertion failed. \n Substring: \n\t\t %s \n Actual: \n\t\t %v", substr, a.actual)
	currentT.Error(errors.StackN(err, 1))
	return false
}

func (a *AnyAssertions) NotEquals(other interface{}) bool {
	if !reflect.DeepEqual(other, a.actual) {
		return true
	}
	err := fmt.Errorf("NotEquals assertion failed. \n Other: \n\t\t [%s] %v \n Actual: \n\t\t [%s] %v", reflect.ValueOf(other).Type(), other, reflect.ValueOf(a.actual).Type(), a.actual)
	currentT.Error(errors.StackN(err, 1))
	return false
}

//IsTrue asserts that actual value is true
func (a *AnyAssertions) IsTrue() bool {
	return a.Equals(true)
}

//IsFalse asserts that actual value is false
func (a *AnyAssertions) IsFalse() bool {
	return a.Equals(false)
}

//IsEmpty asserts that actual value is empty
func (a *AnyAssertions) IsEmpty() bool {
	return a.Equals("")
}

//IsNotEmpty asserts that actual value is not empty
func (a *AnyAssertions) IsNotEmpty() bool {
	return a.NotEquals("")
}

//IsNotNil asserts that actual value is not nil
func (a *AnyAssertions) IsNotNil() bool {
	if a.actual != nil && !reflect.ValueOf(a.actual).IsNil() {
		return true
	}
	err := fmt.Errorf("IsNotNil assertion failed. \n Actual: \n\t\t %v", a.actual)
	currentT.Error(errors.StackN(err, 1))
	return false
}

//IsNil asserts that actual value is nil
func (a *AnyAssertions) IsNil() bool {
	if a.actual == nil || reflect.ValueOf(a.actual).IsNil() {
		return true
	}
	err := fmt.Errorf("IsNil assertion failed. \n Actual: \n\t\t %v", a.actual)
	currentT.Error(errors.StackN(err, 1))
	return false
}

//HasLen asserts that actual value has an expected length
func (a *AnyAssertions) HasLen(expected int) bool {
	length := reflect.ValueOf(a.actual).Len()
	if expected == length {
		return true
	}
	err := fmt.Errorf("HasLen assertion failed. \n Expected: \n\t\t %d \n Actual: \n\t\t %d", expected, length)
	currentT.Error(errors.StackN(err, 1))
	return false
}

//Panics asserts that actual value panics whenever called
func (a *AnyAssertions) Panics() (panicked bool) {
	mustBeFunction(a.actual)
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	reflect.ValueOf(a.actual).Call([]reflect.Value{})
	if !panicked {
		err := errors.New("func.Panics assertion failed. \n Given function didn't panic")
		currentT.Error(errors.StackN(err, 1))
	}
	return
}

//EventuallyEquals asserts that, withtin 5 seconds, the actual function will return same value as expected value
func (a *AnyAssertions) EventuallyEquals(expected interface{}) bool {
	mustBeFunction(a.actual)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		values := reflect.ValueOf(a.actual).Call([]reflect.Value{})
		if reflect.DeepEqual(expected, values[0].Interface()) {
			return true
		}
		select {
		case <-ctx.Done():
			err := fmt.Errorf("func.EventuallyEquals assertion failed. \n Expected: \n\t\t [%s] %s \n Actual: \n\t\t [%s] %s", reflect.ValueOf(expected).Type(), expected, reflect.ValueOf(a.actual).Type(), a.actual)
			currentT.Error(errors.StackN(err, 1))
			return false
		case <-ticker.C:
		}
	}
}

//TemporarilySimilar asserts that actual value is between a range of other time value
func (a *AnyAssertions) TemporarilySimilar(other time.Time, diff time.Duration) bool {
	actual, ok := a.actual.(time.Time)
	if !ok {
		panic("Value is not a time.Time")
	}
	upperBound := other.Add(diff)
	lowerBound := other.Add(diff * -1)
	if actual.After(lowerBound) && actual.Before(upperBound) {
		return true
	}

	err := fmt.Errorf("time.Similar assertion failed. \n Range: \n\t\t %s ~ %s \n Actual: \n\t\t %s", lowerBound, upperBound, actual)
	currentT.Error(errors.StackN(err, 1))
	return false
}
