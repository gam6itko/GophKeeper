package server

import (
	"context"
)

type PrivateDataType int

const (
	TypeLoginPass = iota
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

type PrivateDataDTO struct {
	datatype PrivateDataType
	data     []byte
	meta     map[string]string
}

type IServer interface {
	Registration(ctx context.Context, dto RegistrationDTO) error

	Login(ctx context.Context, dto LoginDTO) error

	Load(ctx context.Context, id uint32) (*PrivateDataDTO, error)
	Store(ctx context.Context, dto PrivateDataDTO) error
}
