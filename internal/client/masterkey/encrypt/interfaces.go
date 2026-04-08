package encrypt

type ICrypt interface {
	Encrypt(key []byte, data []byte) ([]byte, error)
	Decrypt(key []byte, data []byte) ([]byte, error)

	KeyFit(key []byte) []byte
}
