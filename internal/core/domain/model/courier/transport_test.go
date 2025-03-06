package courier

import (
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTransport(t *testing.T) {
	t.Run(
		"Success", func(t *testing.T) {
			transport, err := NewTransport(1, "Test", 2)
			require.NoError(t, err)

			assert.Equal(t, 1, transport.ID())
			assert.Equal(t, "Test", transport.Name())
			assert.Equal(t, 2, transport.Speed())
		},
	)

	t.Run(
		"Invalid ID", func(t *testing.T) {
			_, err := NewTransport(0, "Test", 2)
			assert.Error(t, err)
		},
	)

	t.Run(
		"Invalid Name", func(t *testing.T) {
			_, err := NewTransport(1, "", 2)
			assert.Error(t, err)
		},
	)

	t.Run(
		"Invalid Speed", func(t *testing.T) {
			_, err := NewTransport(1, "Test", 4)
			assert.Error(t, err)
		},
	)
}

func TestTransport_Move(t *testing.T) {
	t.Run(
		"Success", func(t *testing.T) {
			transport, err := NewTransport(1, "Test", 2)
			require.NoError(t, err)

			start, err := kernel.CreateLocation(1, 1)
			require.NoError(t, err)

			end, err := kernel.CreateLocation(2, 2)
			require.NoError(t, err)

			newLocation, err := transport.Move(start, end)
			require.NoError(t, err)

			assert.Equal(t, 2, newLocation.X())
			assert.Equal(t, 2, newLocation.Y())
		},
	)

	t.Run(
		"Start Equals End", func(t *testing.T) {
			transport, err := NewTransport(1, "Test", 2)
			require.NoError(t, err)

			start, err := kernel.CreateLocation(1, 1)
			require.NoError(t, err)

			end, err := kernel.CreateLocation(1, 1)
			require.NoError(t, err)

			newLocation, err := transport.Move(start, end)
			require.NoError(t, err)

			assert.Equal(t, 1, newLocation.X())
			assert.Equal(t, 1, newLocation.Y())
		},
	)

	t.Run(
		"Start Equals End", func(t *testing.T) {
			transport, err := NewTransport(1, "Test", 3)
			require.NoError(t, err)

			start, err := kernel.CreateLocation(5, 5)
			require.NoError(t, err)

			end, err := kernel.CreateLocation(1, 1)
			require.NoError(t, err)

			newLocation, err := transport.Move(start, end)
			require.NoError(t, err)

			assert.Equal(t, 2, newLocation.X())
			assert.Equal(t, 5, newLocation.Y())
		},
	)
}
