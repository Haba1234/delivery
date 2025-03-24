package courier

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var (
	minSpeed = 1
	maxSpeed = 3
)

type TransportID = uuid.UUID

type Transport struct {
	id    TransportID
	name  string
	speed int
}

func NewTransport(id TransportID, name string, speed int) (*Transport, error) {
	if id == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("id")
	}
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if speed < minSpeed || speed > maxSpeed {
		return nil, errs.NewValueIsOutOfRangeError("speed", speed, minSpeed, maxSpeed)
	}

	return &Transport{
		id:    id,
		name:  name,
		speed: speed,
	}, nil
}

func RestoreTransport(id TransportID, name string, speed int) *Transport {
	return &Transport{
		id:    id,
		name:  name,
		speed: speed,
	}
}

func (t *Transport) ID() TransportID {
	return t.id
}

func (t *Transport) Name() string {
	return t.name
}

func (t *Transport) Speed() int {
	return t.speed
}

func (t *Transport) Equals(other *Transport) bool {
	return t.id == other.id
}

func (t *Transport) move(start, end kernel.Location) (kernel.Location, error) {
	if start.Equals(end) {
		return end, nil
	}

	distance := start.DistanceTo(end)
	stepLength := minValue(t.Speed(), distance)

	deltaX := end.X() - start.X()
	deltaY := end.Y() - start.Y()

	stepX := minValue(abs(deltaX), stepLength) * sign(deltaX)
	stepY := minValue(abs(deltaY), stepLength-abs(stepX)) * sign(deltaY)

	newX := start.X() + stepX
	newY := start.Y() + stepY

	return kernel.CreateLocation(newX, newY)
}

// Вспомогательные функции
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func minValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sign(v int) int {
	if v == 0 {
		return 0
	}

	if v > 0 {
		return 1
	}

	return -1
}
