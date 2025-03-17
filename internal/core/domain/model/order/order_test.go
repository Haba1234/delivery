package order

import (
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	location, err := kernel.CreateLocation(5, 5)
	require.NoError(t, err)

	orderID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			newOrder, err := New(orderID, location)
			require.NoError(t, err)

			assert.Equal(t, orderID, newOrder.ID())
			assert.Equal(t, location, newOrder.Location())
			assert.Equal(t, StatusCreated, newOrder.Status())
			assert.Nil(t, newOrder.CourierID())
		},
	)

	t.Run(
		"orderID = nil", func(t *testing.T) {
			newCourier, err := New(uuid.Nil, location)
			assert.Error(t, err)
			require.Nil(t, newCourier)
		},
	)
}

func TestOrder_AssignCourier(t *testing.T) {
	location, err := kernel.CreateLocation(5, 5)
	require.NoError(t, err)

	orderID := uuid.New()
	courierID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			newOrder, err := New(orderID, location)
			require.NoError(t, err)

			executor, err := courier.New("Jon", "Bike", 2, location)
			require.NoError(t, err)

			err = newOrder.Assign(executor)
			require.NoError(t, err)
			assert.Equal(t, StatusAssigned, newOrder.Status())
			assert.NotNil(t, courierID, newOrder.CourierID())

			err = newOrder.Assign(executor)
			assert.ErrorIs(t, err, ErrOrderAlreadyAssigned)
		},
	)

	t.Run(
		"courierID = nil", func(t *testing.T) {
			newOrder, err := New(orderID, location)
			require.NoError(t, err)

			assert.Error(t, newOrder.Assign(nil))
		},
	)
}

func TestOrder_Complete(t *testing.T) {
	location, err := kernel.CreateLocation(5, 5)
	require.NoError(t, err)

	orderID := uuid.New()
	courierID := uuid.New()

	newOrder, err := New(orderID, location)
	require.NoError(t, err)

	executor, err := courier.New("Jon", "Bike", 2, location)
	require.NoError(t, err)

	require.NoError(t, newOrder.Assign(executor))

	assert.Equal(t, StatusAssigned, newOrder.Status())
	assert.NotNil(t, courierID, newOrder.CourierID())

	err = newOrder.Complete()
	require.NoError(t, err)
	assert.Equal(t, StatusCompleted, newOrder.Status())

	err = newOrder.Complete()
	assert.Error(t, err)
}
