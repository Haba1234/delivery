package queries

type GetAllCouriers struct {
	isSet bool
}

func NewGetAllCouriers() (GetAllCouriers, error) {
	return GetAllCouriers{
		isSet: true,
	}, nil
}

func (q GetAllCouriers) IsEmpty() bool {
	return !q.isSet
}
