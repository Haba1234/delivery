package errs

import "fmt"

type VersionIsInvalidError struct {
	msg string
}

func NewVersionIsInvalidError(msg string) VersionIsInvalidError {
	return VersionIsInvalidError{msg: msg}
}

func (v VersionIsInvalidError) Error() string {
	return fmt.Sprintf("version is invalid %s", v.msg)
}
