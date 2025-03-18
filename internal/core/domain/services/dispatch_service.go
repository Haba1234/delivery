package services

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"
)

type DispatchService struct{}

func NewDispatchService() *DispatchService {
	return &DispatchService{}
}

func (*DispatchService) Dispatch(o *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if o == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}

	if len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	orderLocation := o.Location()

	var nearestCourier *courier.Courier
	var minDeliveryTime float64

	for i, c := range couriers {
		deliveryTime, err := c.CalculateTimeToLocation(orderLocation)
		if err != nil {
			return nil, err
		}

		if i == 0 || deliveryTime < minDeliveryTime {
			nearestCourier = c
			minDeliveryTime = deliveryTime
		}
	}

	return nearestCourier, nil
}
