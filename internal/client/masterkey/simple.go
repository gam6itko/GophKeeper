package masterkey

import "errors"

var _ IStorage = (*SimpleStorage)(nil)

// SimpleStorage - простое хранилище для разработки.
type SimpleStorage struct {
	key []byte
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{}
}

func (ths *SimpleStorage) Has() bool {
	return len(ths.key) > 0
}

func (ths *SimpleStorage) Load() ([]byte, error) {
	if !ths.Has() {
		return nil, errors.New("no key stored")
	}

	return ths.key, nil
}

func (ths *SimpleStorage) Store(bytes []byte) error {
	ths.key = bytes
	return nil
}

func (ths *SimpleStorage) Clear() {
	ths.key = nil
}
