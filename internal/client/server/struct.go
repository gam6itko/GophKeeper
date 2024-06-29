package server

import (
	"errors"
	"github.com/gam6itko/goph-keeper/internal/utils/luhn"
	"regexp"
	"strconv"
)

var reCardNumber = regexp.MustCompile("^\\d{16}$")
var reExpires = regexp.MustCompile("^\\d{2}/\\d{2}$")
var reCVV = regexp.MustCompile("^\\d{3}$")

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

type LoginPassDTO struct {
	Login    string
	Password string
}

type TextDTO struct {
	Text string
}

type BinaryDTO struct {
	Data []byte
}
