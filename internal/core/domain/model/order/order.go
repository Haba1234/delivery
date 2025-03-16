package order

import (
	"errors"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var ErrOrderNotAssigned = errors.New("order not assigned")

type (
	ID    = uuid.UUID
	Order struct {
		id        ID
		location  kernel.Location
		status    Status
		courierID courier.ID
	}
)

func New(orderID ID, location kernel.Location) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}

	return &Order{
		id:        orderID,
		location:  location,
		status:    StatusCreated,
		courierID: uuid.Nil,
	}, nil
}

func (o *Order) ID() ID {
	return o.id
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) CourierID() courier.ID {
	return o.courierID
}

func (o *Order) Equals(other *Order) bool {
	return o.id == other.id
}

func (o *Order) Assign(courierID courier.ID) error {
	if courierID == uuid.Nil {
		return errs.NewValueIsRequiredError("courierID")
	}

	o.courierID = courierID
	o.status = StatusAssigned

	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return ErrOrderNotAssigned
	}

	o.status = StatusCompleted

	return nil
}
