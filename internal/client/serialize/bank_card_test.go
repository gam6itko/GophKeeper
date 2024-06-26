package serialize

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBankCardDTO_Validate(t *testing.T) {
	validDTO := BankCardDTO{
		Number:  "4929175965841786",
		Expires: "02/22",
		CVV:     "123",
	}

	t.Run("valid", func(t *testing.T) {
		err := validDTO.Validate()
		require.NoError(t, err)
	})

	t.Run("invalid number", func(t *testing.T) {
		dto := validDTO
		dto.Number = ""
		err := dto.Validate()
		require.Error(t, err)
		require.EqualError(t, err, "invalid Number")
	})

	t.Run("invalid luhn", func(t *testing.T) {
		dto := validDTO
		dto.Number = "0000111100002222"
		err := dto.Validate()
		require.Error(t, err)
		require.EqualError(t, err, "luhn checksum mismatch for Number")
	})

	t.Run("invalid expires", func(t *testing.T) {
		dto := validDTO
		dto.Expires = ""
		err := dto.Validate()
		require.Error(t, err)
		require.EqualError(t, err, "invalid Expires")
	})

	t.Run("invalid cvv", func(t *testing.T) {
		dto := validDTO
		dto.CVV = ""
		err := dto.Validate()
		require.Error(t, err)
		require.EqualError(t, err, "invalid CVV")
	})
}

func TestBankCard_Serialize(t *testing.T) {
	dto := BankCardDTO{
		Number:  "4929175965841786",
		Expires: "02/22",
		CVV:     "123",
	}

	ser := BankCard{}
	b, err := ser.Serialize(dto)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	dto2, err := ser.Deserialize(b)
	require.NoError(t, err)
	require.Equal(t, dto, dto2)
}
