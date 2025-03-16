package courier

import (
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourier(t *testing.T) {
	transport, err := NewTransport(uuid.New(), "Bike", 3)
	require.NoError(t, err)

	location, err := kernel.CreateLocation(5, 5)
	require.NoError(t, err)

	t.Run(
		"Success", func(t *testing.T) {
			newCourier, err := New("Jon", transport, location)
			require.NoError(t, err)

			assert.Equal(t, "Jon", newCourier.ID())
			assert.Equal(t, "Test", newCourier.Name())
			assert.Equal(t, transport, newCourier.Transport())
			assert.Equal(t, location, newCourier.Location())
			assert.Equal(t, StatusFree, newCourier.Status())
		},
	)

	t.Run(
		"Invalid Name", func(t *testing.T) {
			newCourier, err := New("", transport, location)
			assert.Error(t, err)
			require.Nil(t, newCourier)
		},
	)

	t.Run(
		"status busy", func(t *testing.T) {
			newCourier, err := New("Jon", transport, location)
			require.NoError(t, err)
			assert.Equal(t, StatusFree, newCourier.Status())

			newCourier.SetBusy()
			assert.Equal(t, StatusBusy, newCourier.Status())
		},
	)

	t.Run(
		"status free", func(t *testing.T) {
			newCourier, err := New("Jon", transport, location)
			require.NoError(t, err)
			assert.Equal(t, StatusFree, newCourier.Status())

			newCourier.SetBusy()
			assert.Equal(t, StatusBusy, newCourier.Status())

			newCourier.SetFree()
			assert.Equal(t, StatusFree, newCourier.Status())
		},
	)
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {
	transport, err := NewTransport(uuid.New(), "Bike", 2)
	require.NoError(t, err)

	courierLocation, err := kernel.CreateLocation(1, 1)
	require.NoError(t, err)

	courier, err := New("Jon", transport, courierLocation)
	require.NoError(t, err)

	tests := []struct {
		name           string
		targetLocation kernel.Location
		expectedTime   float64
	}{
		{
			name:           "Same location",
			targetLocation: courierLocation,
			expectedTime:   0,
		},
		{
			name: "Directly adjacent location",
			targetLocation: func() kernel.Location {
				loc, err := kernel.CreateLocation(courierLocation.X()+1, courierLocation.Y())
				require.NoError(t, err)
				return loc
			}(),
			expectedTime: 0.5, // Distance = 1, Speed = 2 -> Time = 1 / 2
		},
		{
			name: "Diagonal location",
			targetLocation: func() kernel.Location {
				loc, err := kernel.CreateLocation(courierLocation.X()+3, courierLocation.Y()+4)
				require.NoError(t, err)
				return loc
			}(),
			expectedTime: 3.5, // Distance = 7, Speed = 2 -> Time = 7 / 2
		},
		{
			name: "Far location",
			targetLocation: func() kernel.Location {
				loc, err := kernel.CreateLocation(courierLocation.X()+9, courierLocation.Y()+9)
				require.NoError(t, err)
				return loc
			}(),
			expectedTime: 9.0, // Distance = 18, Speed = 2 -> Time = 18 / 2
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actualTime := courier.CalculateTimeToLocation(tt.targetLocation)
				assert.InEpsilon(t, tt.expectedTime, actualTime, 0.01)
			},
		)
	}
}
