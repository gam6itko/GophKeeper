package masterkey

type SimpleStorage struct {
	key []byte
}

func (ths *SimpleStorage) Has() bool {
	return len(ths.key) > 0
}

func (ths *SimpleStorage) Load() ([]byte, error) {
	return ths.key, nil
}

func (ths *SimpleStorage) Store(bytes []byte) error {
	ths.key = bytes
	return nil
}

func (ths *SimpleStorage) Clear() {
	ths.key = nil
}
