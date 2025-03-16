package courier

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type Courier struct {
	id        ID
	name      string
	transport *Transport
	location  kernel.Location
	status    Status
}

func New(name string, transport *Transport, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	return &Courier{
		id:        uuid.New(),
		name:      name,
		transport: transport,
		location:  location,
		status:    StatusFree,
	}, nil
}

func (c *Courier) ID() ID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Transport() *Transport {
	return c.transport
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) Status() Status {
	return c.status
}

func (c *Courier) Equals(other *Courier) bool {
	return c.id == other.id
}

func (c *Courier) SetBusy() {
	c.status = StatusBusy
}

func (c *Courier) SetFree() {
	c.status = StatusFree
}

// CalculateTimeToLocation calculates the time required for the courier to reach
// the specified location based on transport speed.
func (c *Courier) CalculateTimeToLocation(location kernel.Location) float64 {
	distance := c.location.DistanceTo(location)

	return float64(distance) / float64(c.transport.Speed())
}
