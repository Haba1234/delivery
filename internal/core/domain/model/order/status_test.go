package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StatusShouldReturnCorrectName(t *testing.T) {
	assert.Equal(t, "", StatusEmpty.String())
	assert.Equal(t, "Created", StatusCreated.String())
	assert.Equal(t, "Assigned", StatusAssigned.String())
	assert.Equal(t, "Completed", StatusCompleted.String())
}

func Test_StatusShouldBeEqualWhenAllPropertiesEquals(t *testing.T) {
	assert.True(t, StatusEmpty.Equals(StatusEmpty))
	assert.True(t, StatusCreated.Equals(StatusCreated))
	assert.True(t, StatusAssigned.Equals(StatusAssigned))
	assert.True(t, StatusCompleted.Equals(StatusCompleted))
}

func Test_StatusShouldBeNotEqualWhenAllPropertiesEquals(t *testing.T) {
	assert.False(t, StatusEmpty.Equals(StatusCreated))
	assert.False(t, StatusCreated.Equals(StatusCompleted))
	assert.False(t, StatusCompleted.Equals(StatusAssigned))
}

func Test_StatusShouldBeEmpty(t *testing.T) {
	assert.True(t, StatusEmpty.IsEmpty())
}
