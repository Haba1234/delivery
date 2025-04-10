package eventhandlers

import (
	"context"

	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/errs"
)

type IEventHandler[TNotification any] interface {
	Handle(ctx context.Context, notification TNotification) error
}

type OrderCompletedDomainEventHandler struct {
	orderProducer ports.IOrderProducer
}

func NewOrderCompletedDomainEventHandler(orderProducer ports.IOrderProducer) (
	*OrderCompletedDomainEventHandler, error,
) {
	if orderProducer == nil {
		return nil, errs.NewValueIsRequiredError("orderProducer")
	}

	return &OrderCompletedDomainEventHandler{orderProducer: orderProducer}, nil
}

func (eh *OrderCompletedDomainEventHandler) Handle(_ context.Context, domainEvent order.CompletedDomainEvent) error {
	err := eh.orderProducer.Produce(domainEvent)
	if err != nil {
		return err
	}

	return nil
}
