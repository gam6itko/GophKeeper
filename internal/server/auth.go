package server

import (
	"context"
	"database/sql"
	"github.com/gam6itko/goph-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerImpl struct {
	proto.UnimplementedAuthServer
	db *sql.DB
}

func NewAuthServerImpl(db *sql.DB) *AuthServerImpl {
	return &AuthServerImpl{db: db}
}

func (AuthServerImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (AuthServerImpl) Registration(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Registration not implemented")
}
