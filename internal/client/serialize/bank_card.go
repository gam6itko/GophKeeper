package serialize

import (
	"errors"
	"github.com/gam6itko/goph-keeper/internal/utils/luhn"
	"regexp"
	"strconv"
)

var reCardNumber = regexp.MustCompile("^\\d{16}$")
var reExpires = regexp.MustCompile("^\\d{2}/\\d{2}$")
var reCVV = regexp.MustCompile("^\\d{3}$")

const (
	lenNumber  = 16
	lenExpires = 5
	lenCvv     = 3
	lenTotal   = lenNumber + lenExpires + lenCvv
)

// BankCardDTO структура для передачи на сервер.
// todo для экономии места можно оперировать не строками а числами.
type BankCardDTO struct {
	// Number строка из
	Number string
	// Expires - строка формата 'MM/YY'
	Expires string
	// CVV - строка из 3х цифр.
	CVV string
}

func (dto BankCardDTO) Validate() error {
	// Number validation.
	if len(dto.Number) == 0 || !reCardNumber.MatchString(dto.Number) {
		return errors.New("invalid Number")
	}
	num, err := strconv.ParseUint(dto.Number, 10, 64)
	if err != nil {
		return err
	}
	if !luhn.Validate(num) {
		return errors.New("luhn checksum mismatch for Number")
	}

	//Expires validation.
	if !reExpires.MatchString(dto.Expires) {
		return errors.New("invalid Expires")
	}

	//CVV
	if !reCVV.MatchString(dto.CVV) {
		return errors.New("invalid CVV")
	}

	return nil
}

type BankCard struct{}

func (ths BankCard) Serialize(dto BankCardDTO) ([]byte, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	result := make([]byte, 0, lenTotal)
	result = append(result, []byte(dto.Number)...)
	result = append(result, []byte(dto.Expires)...)
	result = append(result, []byte(dto.CVV)...)

	return result, nil
}

func (ths BankCard) Deserialize(b []byte) (BankCardDTO, error) {
	dto := BankCardDTO{}

	if len(b) != lenTotal {
		return dto, errors.New("invalid length")
	}

	var f, t = 0, lenNumber
	dto.Number = string(b[f:t])
	f, t = lenNumber, lenNumber+lenExpires
	dto.Expires = string(b[f:t])
	f, t = lenNumber+lenExpires, lenNumber+lenExpires+lenCvv
	dto.CVV = string(b[f:t])

	return dto, nil
}
