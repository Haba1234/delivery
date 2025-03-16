package courier

import (
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourier(t *testing.T) {
	location, err := kernel.CreateLocation(5, 5)
	require.NoError(t, err)

	t.Run(
		"Success", func(t *testing.T) {
			newCourier, err := New("Jon", "Bike", 2, location)
			require.NoError(t, err)

			assert.Equal(t, "Jon", newCourier.Name())
			assert.NotNil(t, newCourier.Transport())
			assert.Equal(t, 2, newCourier.Transport().Speed())
			assert.Equal(t, "Bike", newCourier.Transport().Name())
			assert.Equal(t, location, newCourier.Location())
			assert.Equal(t, StatusFree, newCourier.Status())
		},
	)

	t.Run(
		"Invalid Name", func(t *testing.T) {
			newCourier, err := New("", "Bike", 2, location)
			assert.Error(t, err)
			require.Nil(t, newCourier)
		},
	)

	t.Run(
		"status busy", func(t *testing.T) {
			newCourier, err := New("Jon", "Bike", 2, location)
			require.NoError(t, err)
			assert.Equal(t, StatusFree, newCourier.Status())

			err = newCourier.SetBusy()
			require.NoError(t, err)
			assert.Equal(t, StatusBusy, newCourier.Status())

			err = newCourier.SetBusy()
			assert.ErrorIs(t, err, ErrCourierIsBusy)
		},
	)

	t.Run(
		"status free", func(t *testing.T) {
			newCourier, err := New("Jon", "Bike", 2, location)
			require.NoError(t, err)
			assert.Equal(t, StatusFree, newCourier.Status())

			err = newCourier.SetBusy()
			require.NoError(t, err)
			assert.Equal(t, StatusBusy, newCourier.Status())

			err = newCourier.SetFree()
			require.NoError(t, err)
			assert.Equal(t, StatusFree, newCourier.Status())
		},
	)
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {
	courierLocation, err := kernel.CreateLocation(1, 1)
	require.NoError(t, err)

	courier, err := New("Jon", "Bike", 2, courierLocation)
	require.NoError(t, err)

	tests := []struct {
		name           string
		targetLocation kernel.Location
		expectedTime   float64
	}{
		{
			name:           "Same location",
			targetLocation: courierLocation,
			expectedTime:   0.0,
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
				actualTime, err := courier.CalculateTimeToLocation(tt.targetLocation)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTime, actualTime)
			},
		)
	}
}
