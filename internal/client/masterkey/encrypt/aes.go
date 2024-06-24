package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type AESCrypt struct {
}

func NewAESCrypt() *AESCrypt {
	return &AESCrypt{}
}

func (ths AESCrypt) Encrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	result := make([]byte, aes.BlockSize+len(data))
	iv := result[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(result[aes.BlockSize:], data)

	return result, nil
}

func (ths AESCrypt) Decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("invalid ciphertext block size")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}

func (ths AESCrypt) KeyFit(key []byte) []byte {
	keyLen := len(key)
	noFit := keyLen % aes.BlockSize
	if noFit == 0 {
		return key
	}

	c := ((keyLen / aes.BlockSize) + 1) * aes.BlockSize
	newKey := make([]byte, keyLen, c)
	copy(newKey, key)
	addLen := aes.BlockSize - noFit
	for i := 0; i < addLen; i++ {
		newKey = append(newKey, 0)
	}
	return newKey
}
