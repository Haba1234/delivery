package order

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
)

func toModelDB(aggregate *order.Order) ModelOrder {
	return ModelOrder{
		ID:        aggregate.ID(),
		CourierID: aggregate.CourierID(),
		Location: ModelLocation{
			X: aggregate.Location().X(),
			Y: aggregate.Location().Y(),
		},
		Status: aggregate.Status(),
	}
}

func toDomain(model ModelOrder) *order.Order {
	return order.Restore(
		model.ID,
		model.CourierID,
		kernel.RestoreLocation(model.Location.X, model.Location.Y),
		model.Status,
	)
}
