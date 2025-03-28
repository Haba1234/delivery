package services

import (
	"errors"
	"math"

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
	minDeliveryTime := float64(math.MaxInt64)

	for i, c := range couriers {
		if c == nil || c.Status() == courier.StatusBusy {
			continue // Пропускаем nil курьеров или уже занятых курьеров
		}

		deliveryTime, err := c.CalculateTimeToLocation(orderLocation)
		if err != nil {
			return nil, err
		}

		if i == 0 || deliveryTime < minDeliveryTime {
			nearestCourier = c
			minDeliveryTime = deliveryTime
		}
	}

	if nearestCourier == nil {
		return nil, errors.New("no available couriers found")
	}

	err := o.Assign(nearestCourier)
	if err != nil {
		return nil, err
	}

	err = nearestCourier.SetBusy()
	if err != nil {
		return nil, err
	}

	return nearestCourier, nil
}
