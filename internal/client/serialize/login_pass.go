package serialize

import (
	"errors"
	"log"
)

type LoginPassDTO struct {
	Login    string
	Password string
}

// LoginPass преобразует структуру данных в массив.
// Разделителем является 0-byte.
type LoginPass struct{}

func (ths LoginPass) Serialize(dto LoginPassDTO) ([]byte, error) {
	if dto.Login == "" || dto.Password == "" {
		return nil, errors.New("empty login or password")
	}

	result := make([]byte, 0, len(dto.Login)+len(dto.Password)+1)
	result = append(result, []byte(dto.Login)...)
	result = append(result, 0)
	result = append(result, []byte(dto.Password)...)

	return result, nil
}

func (ths LoginPass) Deserialize(data []byte) (LoginPassDTO, error) {
	z, ok := findZeroByteIndex(data, 0)
	if !ok {
		log.Fatal("invalid format")
	}
	return LoginPassDTO{
		Login:    string(data[:z]),
		Password: string(data[z+1:]),
	}, nil
}

// findZeroByteIndex возвращает индекс первого попавшегося элемента со значением 0.
// Поиска начинается с индекса from.
func findZeroByteIndex(data []byte, from int) (int, bool) {
	l := len(data)
	for i := from; i < l; i++ {
		if data[i] == 0 {
			return i, true
		}
	}

	return 0, false
}
