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

func caller() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	file = trimBasePath(file)
	return file, line
}

//New creates a new error
func New(text string) error {
	file, line := caller()
	return stdError.New(fmt.Sprintf("%s (%s:%d)", text, file, line))
}

//Wrap existing error with additional text
func Wrap(err error, format string, a ...interface{}) error {
	text := fmt.Sprintf(format, a...)
	file, line := caller()

	casted, ok := err.(*theError)
	if ok {
		text = fmt.Sprintf("%s (%s:%d)\n- %s", text, file, line, casted.text)
	} else {
		text = fmt.Sprintf("%s (%s:%d)", text, file, line)
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
