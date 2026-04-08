package luhn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	assert.True(t, Validate(4561_2612_1234_5467))
	assert.True(t, Validate(5580_4733_7202_4733))

	assert.False(t, Validate(5580_4733_7202_4732))
}
