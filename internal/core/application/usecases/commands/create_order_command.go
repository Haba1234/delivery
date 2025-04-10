package commands

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrder struct {
	orderID order.OrderID
	street  string

	isSet bool
}

func (c CreateOrder) OrderID() order.OrderID {
	return c.orderID
}

func (c CreateOrder) Street() string {
	return c.street
}

func NewCreateOrder(orderID order.OrderID, street string) (CreateOrder, error) {
	if orderID == uuid.Nil {
		return CreateOrder{}, errs.NewValueIsInvalidError("orderID")
	}
	if street == "" {
		return CreateOrder{}, errs.NewValueIsRequiredError("street")
	}

	return CreateOrder{
		orderID: orderID,
		street:  street,

		isSet: true,
	}, nil
}

func (c CreateOrder) IsEmpty() bool {
	return !c.isSet
}
