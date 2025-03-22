package courier

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
)

func toModelDB(aggregate *courier.Courier) ModelCourier {
	return ModelCourier{
		ID:   aggregate.ID(),
		Name: aggregate.Name(),
		Transport: ModelTransport{
			ID:        aggregate.Transport().ID(),
			Name:      aggregate.Transport().Name(),
			Speed:     aggregate.Transport().Speed(),
			CourierID: aggregate.ID(),
		},
		Location: ModelLocation{
			X: aggregate.Location().X(),
			Y: aggregate.Location().Y(),
		},
		Status: aggregate.Status(),
	}
}

func toDomain(model ModelCourier) *courier.Courier {
	return courier.Restore(
		model.ID,
		model.Name,
		courier.RestoreTransport(model.Transport.ID, model.Transport.Name, model.Transport.Speed),
		kernel.RestoreLocation(model.Location.X, model.Location.Y),
		model.Status,
	)
}
