package queries

type GetNotCompletedOrders struct {
	isSet bool
}

func NewGetNotCompletedOrders() (GetNotCompletedOrders, error) {
	return GetNotCompletedOrders{
		isSet: true,
	}, nil
}

func (q GetNotCompletedOrders) IsEmpty() bool {
	return !q.isSet
}
