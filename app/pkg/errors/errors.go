package errors

import (
	stdError "errors"
	"fmt"
	"runtime"
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
func New(format string, a ...interface{}) error {
	file, line := caller(2)
	text := fmt.Sprintf(format, a...)
	return stdError.New(fmt.Sprintf("%s (%s:%d)", text, file, line))
}

//Wrap existing error with additional text
func Wrap(err error, format string, a ...interface{}) error {
	return wrap(err, 0, format, a...)
}

//Stack add current code location without adding more info
func Stack(err error) error {
	return wrap(err, 0, "")
}

//StackN add code location of N steps before without adding more info
func StackN(err error, skip int) error {
	return wrap(err, skip, "")
}

func wrap(err error, skip int, format string, a ...interface{}) error {
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
