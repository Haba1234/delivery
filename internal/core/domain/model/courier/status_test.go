package courier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StatusShouldReturnCorrectName(t *testing.T) {
	assert.Equal(t, "", StatusEmpty.String())
	assert.Equal(t, "free", StatusFree.String())
	assert.Equal(t, "busy", StatusBusy.String())
}

func Test_StatusShouldBeEqualWhenAllPropertiesEquals(t *testing.T) {
	assert.True(t, StatusEmpty.Equals(StatusEmpty))
	assert.True(t, StatusFree.Equals(StatusFree))
	assert.True(t, StatusBusy.Equals(StatusBusy))
}

func Test_StatusShouldBeNotEqualWhenAllPropertiesEquals(t *testing.T) {
	assert.False(t, StatusEmpty.Equals(StatusFree))
	assert.False(t, StatusBusy.Equals(StatusFree))
}

func Test_StatusShouldBeEmpty(t *testing.T) {
	assert.True(t, StatusEmpty.IsEmpty())
}
