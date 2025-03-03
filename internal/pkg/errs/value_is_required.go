package errs

import "fmt"

type ValueIsRequiredError struct {
	msg string
}

func NewValueIsRequiredError(msg string) ValueIsRequiredError {
	return ValueIsRequiredError{msg: msg}
}

func (v ValueIsRequiredError) Error() string {
	return fmt.Sprintf("value is required %s", v.msg)
}
