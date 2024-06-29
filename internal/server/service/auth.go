package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gam6itko/goph-keeper/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type AuthServerImpl struct {
	proto.UnimplementedAuthServer
	db *sql.DB
}

func NewAuthServerImpl(db *sql.DB) *AuthServerImpl {
	return &AuthServerImpl{db: db}
}

// Registration - регистрация нового пользователя.
func (ths AuthServerImpl) Registration(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	if len(req.Username) < 3 {
		return nil, status.Error(codes.InvalidArgument, "username is too short")
	}
	if len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password is too short")
	}

	row := ths.db.QueryRowContext(
		ctx,
		"SELECT COUNT(1) FROM `user` WHERE `username` = ?",
		req.Username,
	)
	if row.Err() != nil {
		log.Printf("ERROR. %s", row.Err())
		errMessage := fmt.Sprintf("query row error: %s", row.Err())
		return nil, status.Error(codes.Internal, errMessage)
	}

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		errMessage := fmt.Sprintf("scan row error: %s", row.Err())
		log.Printf("ERROR. %s", errMessage)
		return nil, status.Error(codes.Internal, errMessage)
	}
	if cnt > 0 {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	b, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errMessage := fmt.Sprintf("password hash error: %s", row.Err())
		log.Printf("ERROR. %s", errMessage)
		return nil, status.Error(codes.Internal, errMessage)
	}

	_, err = ths.db.ExecContext(
		ctx,
		"INSERT INTO `user` (`username`, `password`) VALUES (?, ?)`",
		req.Username,
		b,
	)
	if err != nil {
		errMessage := fmt.Sprintf("registration user save error: %s", row.Err())
		log.Printf("ERROR. %s", errMessage)
		return nil, status.Error(codes.Internal, errMessage)
	}

	return &proto.RegistrationResponse{}, nil
}

// Login - аутентификация пользователя. Создаёт новый JWT токен.
func (AuthServerImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if len(req.Username) < 3 {
		return nil, status.Error(codes.InvalidArgument, "username is too short")
	}
	if len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password is too short")
	}

	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
