package kernel

import (
	"errors"
	"math/rand"
)

var (
	ErrLocationValueIsRequired = errors.New("coordinates outside the permitted range")
)

type Location struct {
	x int
	y int
}

func CreateLocation(x int, y int) (Location, error) {
	if x < 1 || y < 1 || x > 10 || y > 10 {
		return Location{}, ErrLocationValueIsRequired
	}

	return Location{x, y}, nil
}

func CreateRandomLocation() Location {
	return Location{rand.Intn(10) + 1, rand.Intn(10) + 1}
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equals(l2 Location) bool {
	return l == l2
}

func (l Location) DistanceTo(l2 Location) int {
	return abs(l.x-l2.x) + abs(l.y-l2.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
