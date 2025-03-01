package errs

import "fmt"

type ValueIsInvalidError struct {
	msg string
}

func NewValueIsInvalidError(msg string) ValueIsInvalidError {
	return ValueIsInvalidError{msg: msg}
}

func (v ValueIsInvalidError) Error() string {
	return fmt.Sprintf("value is invalid %s", v.msg)
}
