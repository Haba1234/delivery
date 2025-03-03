package kernel

import (
	"errors"
	"math/rand"
)

var (
	ErrLocationValueIsRequired = errors.New("coordinates outside the permitted range")
)

var (
	minCoordinate = 1
	maxCoordinate = 10
)

type Location struct {
	x int
	y int
}

func CreateLocation(x int, y int) (Location, error) {
	if x < minCoordinate || y < minCoordinate || x > maxCoordinate || y > maxCoordinate {
		return Location{}, ErrLocationValueIsRequired
	}

	return Location{x, y}, nil
}

func CreateRandomLocation() (Location, error) {
	x := rand.Intn(maxCoordinate) + 1
	y := rand.Intn(maxCoordinate) + 1
	return CreateLocation(x, y)
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
