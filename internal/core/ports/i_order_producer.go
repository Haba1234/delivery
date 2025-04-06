package ports

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
)

type IOrderProducer interface {
	Produce(event order.CompletedDomainEvent) error
	Close() error
}
