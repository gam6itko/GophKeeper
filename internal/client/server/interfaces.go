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

type ILoginServer interface {
	Login(ctx context.Context, dto LoginDTO) error
}

type IRegistrationServer interface {
	Register(ctx context.Context, dto RegistrationDTO) error
}

type IPrivateStorage interface {
	Load(ctx context.Context, id uint32) (PrivateDataDTO, error)
	Store(ctx context.Context, dto PrivateDataDTO) error
}
