package server

import (
	"context"
	"fmt"
)

// TPrivateData - тип данных хранящихся на сервере.
type TPrivateData int

func (ths TPrivateData) String() string {
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
	TypeUndefined TPrivateData = iota
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
	Type TPrivateData
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

	// Login - аутентификация.
	Login(ctx context.Context, dto LoginDTO) error
	// Logout - сбросить аутентификацию.
	Logout(ctx context.Context) error

	// List - получить список данных пользователя на сервере.
	List(ctx context.Context) ([]PrivateDataListItemDTO, error)
	// Load - загрузить одну запись с данными пользователя.
	Load(ctx context.Context, id uint32) (*PrivateDataDTO, error)
	// Store - сохранить одну запись с данными пользователя.
	Store(ctx context.Context, dto PrivateDataDTO) error
}
