package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

type theError struct {
	text  string
	cause error
}

func caller(skip int) (string, int) {
	_, file, line, _ := runtime.Caller(skip)
	file = trimBasePath(file)
	return file, line
}

//New creates a new error
func New(format string, a ...any) error {
	file, line := caller(2)
	text := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s (%s:%d)", text, file, line)
}

//Wrap existing error with additional text
func Wrap(err error, format string, a ...any) error {
	return wrap(err, 0, format, a...)
}

//Panicked wraps panick errow with extra details of error
func Panicked(r any) error {
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}
	err = Wrap(err, identifyPanic())
	return Wrap(err, string(debug.Stack()))
}

//Stack add current code location without adding more info
func Stack(err error) error {
	return wrap(err, 0, "")
}

//StackN add code location of N steps before without adding more info
func StackN(err error, skip int) error {
	return wrap(err, skip, "")
}

func wrap(err error, skip int, format string, a ...any) error {
	if err == nil {
		return err
	}

	file, line := caller(3 + skip)

	text := ""
	if format != "" {
		text = fmt.Sprintf("%s (%s:%d)", fmt.Sprintf(format, a...), file, line)
	} else {
		text = fmt.Sprintf("%s:%d", file, line)
	}

	casted, ok := err.(*theError)
	if ok {
		text = fmt.Sprintf("%s\n- %s", text, casted.text)
	}

	return &theError{
		text:  text,
		cause: Cause(err),
	}
}

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(4, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}

//Cause returns the original cause of the error stack
func Cause(err error) error {
	casted, ok := err.(*theError)
	if ok {
		return casted.cause
	}
	return err
}

//Error formats the error trace into human readable
func (err *theError) Error() string {
	return fmt.Sprintf("Error Trace: \n- %s\n- %v", err.text, err.cause)
}
