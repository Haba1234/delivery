package commands

type MoveCouriers struct {
	isSet bool
}

func NewMoveCouriers() (MoveCouriers, error) {
	return MoveCouriers{
		isSet: true,
	}, nil
}

func (c MoveCouriers) IsEmpty() bool {
	return !c.isSet
}
