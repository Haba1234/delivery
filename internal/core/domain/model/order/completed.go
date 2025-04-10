package order

import "github.com/google/uuid"

const CompletedDomainEventName = "order.completed.event"

type DomainEventID = uuid.UUID

type CompletedDomainEvent struct {
	id   DomainEventID
	name string

	orderID     OrderID
	orderStatus string

	isSet bool
}

func NewCompletedDomainEvent(aggregate *Order) CompletedDomainEvent {
	return CompletedDomainEvent{
		id:   uuid.New(),
		name: CompletedDomainEventName,

		orderID:     aggregate.ID(),
		orderStatus: aggregate.Status().String(),

		isSet: true,
	}
}

func (e CompletedDomainEvent) IsEmpty() bool {
	return !e.isSet
}

func (e CompletedDomainEvent) ID() DomainEventID { return e.id }

func (e CompletedDomainEvent) Name() string {
	return e.name
}

func (e CompletedDomainEvent) OrderID() OrderID {
	return e.orderID
}

func (e CompletedDomainEvent) OrderStatus() string {
	return e.orderStatus
}
