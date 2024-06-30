package encrypt

import (
	"crypto/aes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAESCrypt_Encrypt(t *testing.T) {
	t.Run("key size match", func(t *testing.T) {
		key := []byte("0123456789abcdef")
		data := []byte("Lorem ipsum dolor sit amet")

		c := AESCrypt{}

		encrypted, err := c.Encrypt(key, data)
		require.Nil(t, err)

		decrypted, err := c.Decrypt(key, encrypted)
		require.Nil(t, err)
		require.Equal(t, data, decrypted)
	})

	t.Run("key size not match", func(t *testing.T) {
		c := AESCrypt{}

		key := c.KeyFit([]byte("012345"))
		data := []byte("Lorem ipsum dolor sit amet")

		encrypted, err := c.Encrypt(key, data)
		require.Nil(t, err)

		decrypted, err := c.Decrypt(key, encrypted)
		require.Nil(t, err)
		require.Equal(t, data, decrypted)
	})
}

func TestAESCrypt_KeyFit(t *testing.T) {
	c := AESCrypt{}
	key := c.KeyFit([]byte("0"))
	require.Len(t, key, aes.BlockSize)
}
