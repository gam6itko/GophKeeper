package serialize

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginPass_Serialize(t *testing.T) {
	dto := LoginPassDTO{
		Login:    "username1",
		Password: "password2",
	}
	e := LoginPass{}
	b, err := e.Serialize(dto)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	dto2, err := e.Deserialize(b)
	require.NoError(t, err)
	assert.Equal(t, dto.Login, dto2.Login)
	assert.Equal(t, dto.Password, dto2.Password)
}
