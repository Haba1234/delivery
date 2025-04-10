package order

import (
	"errors"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/pkg/ddd"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var (
	ErrOrderNotAssigned     = errors.New("order not assigned")
	ErrThisCourierIsBusy    = errors.New("this courier is already busy")
	ErrOrderAlreadyAssigned = errors.New("order has already been appointed courier")
)

type (
	OrderID = uuid.UUID
	Order   struct {
		id        OrderID
		courierID *courier.CourierID
		location  kernel.Location
		status    Status

		domainEvents []ddd.IDomainEvent
	}
)

func New(orderID OrderID, location kernel.Location) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}

	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}

	return &Order{
		id:       orderID,
		location: location,
		status:   StatusCreated,
	}, nil
}

func Restore(orderID OrderID, courierID *courier.CourierID, location kernel.Location, status Status) *Order {
	return &Order{
		id:        orderID,
		location:  location,
		status:    status,
		courierID: courierID,
	}
}

func (o *Order) ID() OrderID {
	return o.id
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) CourierID() *courier.CourierID {
	return o.courierID
}

func (o *Order) Equals(other *Order) bool {
	return o.id == other.id
}

func (o *Order) Assign(executor *courier.Courier) error {
	if executor == nil {
		return errs.NewValueIsRequiredError("executor")
	}

	if executor.Status() == courier.StatusBusy {
		return ErrThisCourierIsBusy
	}

	if o.status != StatusCreated {
		return ErrOrderAlreadyAssigned
	}

	id := executor.ID()
	o.courierID = &id
	o.status = StatusAssigned

	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return ErrOrderNotAssigned
	}

	o.status = StatusCompleted

	// Опубликовать доменное событие
	o.raiseDomainEvent(NewCompletedDomainEvent(o))

	return nil
}

func (o *Order) IsCompleted() bool {
	return o.status == StatusCompleted
}

func (o *Order) ClearDomainEvents() {
	o.domainEvents = []ddd.IDomainEvent{}
}

func (o *Order) GetDomainEvents() []ddd.IDomainEvent {
	return o.domainEvents
}

func (o *Order) raiseDomainEvent(event ddd.IDomainEvent) {
	o.domainEvents = append(o.domainEvents, event)
}
