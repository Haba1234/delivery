package courier

import (
	"errors"

	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var ErrCourierIsBusy = errors.New("courier is busy")

type ID = uuid.UUID

type Courier struct {
	id        ID
	name      string
	transport *Transport
	location  kernel.Location
	status    Status
}

func New(name, transportName string, transportSpeed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if transportName == "" {
		return nil, errs.NewValueIsRequiredError("transportName")
	}

	if transportSpeed <= 0 {
		return nil, errs.NewValueIsRequiredError("transportSpeed")
	}

	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}

	transport, err := NewTransport(uuid.New(), transportName, transportSpeed)
	if err != nil {
		return nil, err
	}

	return &Courier{
		id:        uuid.New(),
		name:      name,
		transport: transport,
		location:  location,
		status:    StatusFree,
	}, nil
}

func Restore(id ID, name string, transport *Transport, location kernel.Location, status Status) *Courier {
	return &Courier{
		id:        id,
		name:      name,
		transport: transport,
		location:  location,
		status:    status,
	}
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

func (c *Courier) SetBusy() error {
	if c.status == StatusBusy {
		return ErrCourierIsBusy
	}
	c.status = StatusBusy
	return nil
}

func (c *Courier) SetFree() error {
	c.status = StatusFree
	return nil
}

// CalculateTimeToLocation calculates the time required for the courier to reach
// the specified location based on transport speed.
func (c *Courier) CalculateTimeToLocation(location kernel.Location) (float64, error) {
	if location.IsEmpty() {
		return 0, errs.NewValueIsRequiredError("location")
	}

	distance := c.location.DistanceTo(location)

	return float64(distance) / float64(c.transport.Speed()), nil
}
