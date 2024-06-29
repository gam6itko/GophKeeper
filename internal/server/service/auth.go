package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gam6itko/goph-keeper/internal/server/jwt"
	"github.com/gam6itko/goph-keeper/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type AuthServerImpl struct {
	proto.UnimplementedAuthServer
	db     *sql.DB
	issuer *jwt.Issuer
}

func NewAuthServerImpl(db *sql.DB, issuer *jwt.Issuer) *AuthServerImpl {
	return &AuthServerImpl{
		db:     db,
		issuer: issuer,
	}
}

// Login - аутентификация пользователя. Создаёт новый JWT токен.
func (ths AuthServerImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if len(req.Username) < 3 {
		return nil, status.Error(codes.InvalidArgument, "username is too short")
	}
	if len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password is too short")
	}

	row := ths.db.QueryRowContext(
		ctx,
		"SELECT `id`, `password` FROM `user` WHERE `username` = ?",
		req.Username,
	)
	if row.Err() != nil {
		log.Printf("ERROR. %s", row.Err())
		errMessage := fmt.Sprintf("query row error: %s", row.Err())
		return nil, status.Error(codes.Internal, errMessage)
	}

	var (
		userID uint64
		pass   []byte
	)
	if err := row.Scan(&userID, &pass); err != nil {
		log.Printf("ERROR. %s", err)
		return nil, status.Error(codes.Unauthenticated, "username or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword(pass, []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "username or password is incorrect")
	}

	token, err := ths.issuer.Issue(userID)
	if err != nil {
		log.Printf("ERROR. %s", err)
		return nil, status.Error(codes.Internal, "failed to issue token")
	}

	return &proto.LoginResponse{Token: token}, nil
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
		"INSERT INTO `user` (`username`, `password`) VALUES (?, ?)",
		req.Username,
		b,
	)
	if err != nil {
		log.Printf("ERROR. %s", fmt.Sprintf("registration user save error: %s", row.Err()))
		return nil, status.Error(codes.Internal, "registration user save error")
	}

	return &proto.RegistrationResponse{}, nil
}
