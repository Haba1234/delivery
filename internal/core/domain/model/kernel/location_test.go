package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocation_CreateLocation(t *testing.T) {
	tests := []struct {
		name    string
		x, y    int
		wantErr error
	}{
		{"valid coordinates", 1, 2, nil},
		{"negative x", -1, 0, ErrLocationValueIsRequired},
		{"negative y", 0, -1, ErrLocationValueIsRequired},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				l, err := CreateLocation(tt.x, tt.y)
				require.ErrorIs(t, err, tt.wantErr)

				if nil == err {
					assert.Equal(t, l.X(), tt.x)
					assert.Equal(t, l.Y(), tt.y)
				}
			},
		)
	}
}

func TestLocation_Equals(t *testing.T) {
	tests := []struct {
		name   string
		l1, l2 Location
		want   bool
	}{
		{"same locations", Location{1, 2}, Location{1, 2}, true},
		{"different x", Location{1, 2}, Location{3, 2}, false},
		{"different y", Location{1, 2}, Location{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, tt.l1.Equals(tt.l2))
			},
		)
	}
}

func TestLocation_DistanceTo(t *testing.T) {
	tests := []struct {
		name string
		l1   Location
		l2   Location
		want int
	}{
		{
			name: "test 1",
			l1:   Location{4, 9},
			l2:   Location{2, 6},
			want: 5,
		},
		{
			name: "test 2",
			l1:   Location{1, 1},
			l2:   Location{5, 1},
			want: 4,
		},
		{
			name: "test 3",
			l1:   Location{5, 5},
			l2:   Location{10, 10},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, tt.l1.DistanceTo(tt.l2))
			},
		)
	}
}
