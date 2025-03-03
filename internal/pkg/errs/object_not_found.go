package errs

import "fmt"

type ObjectNotFoundError struct {
	msg string
}

func NewObjectNotFoundError(msg string) ObjectNotFoundError {
	return ObjectNotFoundError{msg: msg}
}

func (e ObjectNotFoundError) Error() string {
	return fmt.Sprintf("object not found %s", e.msg)
}
