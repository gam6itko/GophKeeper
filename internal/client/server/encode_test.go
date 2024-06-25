package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginPassEncoder(t *testing.T) {
	dto := LoginPassDTO{
		Login:    "username1",
		Password: "password2",
	}
	e := LoginPassEncoder{}
	b, err := e.Encode(dto)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	dto2, err := e.Decode(b)
	require.NoError(t, err)
	assert.Equal(t, dto.Login, dto2.Login)
	assert.Equal(t, dto.Password, dto2.Password)
}
