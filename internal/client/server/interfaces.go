package server

import (
	"context"
	"fmt"
)

type PrivateDataType int

func (ths PrivateDataType) String() string {
	switch ths {
	case TypeLoginPass:
		return "login,pass"
	case TypeText:
		return "text"
	case TypeBinary:
		return "binary"
	case TypeBankCard:
		return "bank card"
	default:
		return fmt.Sprintf("type: %d", int(ths))
	}
}

const (
	TypeUndefined = iota
	TypeLoginPass
	TypeText
	TypeBinary
	TypeBankCard
)

type LoginDTO struct {
	Username string
	Password string
}

type RegistrationDTO struct {
	Username string
	Password string
}

// BasePrivateDataDTO общая часть.
type BasePrivateDataDTO struct {
	ID   uint32
	Name string
	Type PrivateDataType
	Meta string
}

// PrivateDataDTO данные пришли с сервера.
// Так же используется для сохранения.
type PrivateDataDTO struct {
	BasePrivateDataDTO
	Data []byte
}

// PrivateDataListItemDTO строчка в списке данных на сервере.
type PrivateDataListItemDTO struct {
	BasePrivateDataDTO
}

type IServer interface {
	Registration(ctx context.Context, dto RegistrationDTO) error

	Login(ctx context.Context, dto LoginDTO) error
	Logout(ctx context.Context) error

	List(ctx context.Context) ([]PrivateDataListItemDTO, error)
	Load(ctx context.Context, id uint32) (*PrivateDataDTO, error)
	Store(ctx context.Context, dto PrivateDataDTO) error
}
