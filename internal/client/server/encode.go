package server

import (
	"errors"
	"log"
)

type LoginPassDTO struct {
	Login    string
	Password string
}

// LoginPassEncoder преобразует структуру данных в массив.
// Разделителем является 0-byte.
type LoginPassEncoder struct{}

func (e LoginPassEncoder) Encode(dto LoginPassDTO) ([]byte, error) {
	if dto.Login == "" || dto.Password == "" {
		return nil, errors.New("empty login or password")
	}

	result := make([]byte, 0, len(dto.Login)+len(dto.Password)+1)
	result = append(result, []byte(dto.Login)...)
	result = append(result, 0)
	result = append(result, []byte(dto.Password)...)

	return result, nil
}

func (e LoginPassEncoder) Decode(data []byte) (LoginPassDTO, error) {
	z, ok := findZeroByteIndex(data, 0)
	if !ok {
		log.Fatal("invalid format")
	}
	return LoginPassDTO{
		Login:    string(data[:z]),
		Password: string(data[z+1:]),
	}, nil
}

type TextEncoder struct{}

func (e TextEncoder) Encode(text string) ([]byte, error) {
	return []byte(text), nil
}

func (e TextEncoder) Decode(data []byte) (string, error) {
	return string(data), nil
}

func findZeroByteIndex(data []byte, from int) (int, bool) {
	l := len(data)
	for i := 0; i < l; i++ {
		if data[i] == 0 {
			return i, true
		}
	}

	return 0, false
}
