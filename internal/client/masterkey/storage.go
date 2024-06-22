package masterkey

import (
	"errors"
	"github.com/awnumar/memguard"
)

type IStorage interface {
	Has() bool
	Load() ([]byte, error)
	Store([]byte) error
}

type MemGuardStorage struct {
	enclave *memguard.Enclave
}

func NewMemGuardStorage() *MemGuardStorage {
	return &MemGuardStorage{}
}

func (ths *MemGuardStorage) Has() bool {
	if ths.enclave == nil {
		return false
	} else {
		return true
	}
}

func (ths *MemGuardStorage) Load() ([]byte, error) {
	if ths.enclave == nil {
		return nil, errors.New("enclave not initialized")
	}

	b, err := ths.enclave.Open()
	if err != nil {
		return nil, err
	}
	defer b.Destroy()

	return b.Bytes(), nil
}

func (ths *MemGuardStorage) Store(key []byte) error {
	ths.enclave = memguard.NewEnclave(key)

	return nil
}
