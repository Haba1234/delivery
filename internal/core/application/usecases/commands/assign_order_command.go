package commands

type AssignOrder struct {
	isSet bool
}

func NewAssignOrder() (AssignOrder, error) {
	return AssignOrder{
		isSet: true,
	}, nil
}

func (c AssignOrder) IsEmpty() bool {
	return !c.isSet
}
