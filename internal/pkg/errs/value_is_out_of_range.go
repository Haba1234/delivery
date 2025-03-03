package errs

import "fmt"

func NewValueIsOutOfRangeError(paramName string, value any, min any, max any) ValueIsOutOfRange {
	return ValueIsOutOfRange{
		paramName: paramName,
		value:     value,
		min:       min,
		max:       max,
	}
}

type ValueIsOutOfRange struct {
	paramName string
	value     any
	min       any
	max       any
}

func (v ValueIsOutOfRange) Error() string {
	return fmt.Sprintf(
		"value %v of %v is out of range , min value is %v, max value is %v",
		v.value, v.paramName, v.min, v.max,
	)
}
