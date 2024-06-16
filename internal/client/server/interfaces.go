package server

import (
	"context"
)

type ILoginServer interface {
	Login(ctx context.Context, username string, password string) error
}

type IRegistrationServer interface {
	Register(ctx context.Context, username string, password string) error
}
