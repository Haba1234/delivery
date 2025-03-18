package services

import (
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDispatchService_Dispatch(t *testing.T) {
	ds := NewDispatchService()

	// create test order
	orderID := uuid.New()

	orderLocation, err := kernel.CreateLocation(5, 5)
	require.NoError(t, err)

	testOrder, err := order.New(orderID, orderLocation)
	require.NoError(t, err)

	// create test couriers

	t.Run(
		"nil order", func(t *testing.T) {
			var couriers []*courier.Courier
			result, err := ds.Dispatch(nil, couriers)
			assert.Nil(t, result)
			assert.IsType(t, errs.NewValueIsRequiredError(""), err)
			assert.Equal(t, "value is required order", err.Error())
		},
	)

	t.Run(
		"empty courier list", func(t *testing.T) {
			result, err := ds.Dispatch(testOrder, []*courier.Courier{})
			assert.Nil(t, result)
			assert.IsType(t, errs.NewValueIsRequiredError(""), err)
			assert.Equal(t, "value is required couriers", err.Error())
		},
	)

	t.Run(
		"dispatch to nearest courier", func(t *testing.T) {
			// Создаем курьеров с фиксированными координатами
			c1Location, err := kernel.CreateLocation(1, 1) // Расстояние 8
			require.NoError(t, err)
			c2Location, err := kernel.CreateLocation(4, 4) // Расстояние 2
			require.NoError(t, err)
			c3Location, err := kernel.CreateLocation(7, 7) // Расстояние 4
			require.NoError(t, err)

			c1, err := courier.New("courier 1", "name 1", 1, c1Location)
			require.NoError(t, err)
			c2, err := courier.New("courier 2", "name 2", 2, c2Location)
			require.NoError(t, err)
			c3, err := courier.New("courier 3", "name 3", 3, c3Location)
			require.NoError(t, err)

			// Проверяем, что Dispatch вернет курьера с ближайшим расстоянием
			result, err := ds.Dispatch(testOrder, []*courier.Courier{c1, c2, c3})
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, c2, result) // Ближайший курьер c2
		},
	)
}
