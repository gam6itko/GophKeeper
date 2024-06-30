package service

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gam6itko/goph-keeper/internal/server/jwt"
	"github.com/gam6itko/goph-keeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthServerImpl_Registration(t *testing.T) {
	t.Run("username is too short", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.RegistrationRequest{}
		resp, err := svc.Registration(context.TODO(), &req)
		require.Nil(t, resp)
		require.Error(t, err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = username is too short")
	})

	t.Run("password is too short", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.RegistrationRequest{
			Username: "gam6itko@gmail.com",
		}
		resp, err := svc.Registration(context.TODO(), &req)
		require.Nil(t, resp)
		require.Error(t, err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = password is too short")
	})

	t.Run("user already exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		// Пользователь с таким username уже есть.
		mock.
			ExpectQuery("SELECT (.+) FROM `user`").
			WillReturnRows(
				sqlmock.NewRows([]string{"cnt"}).
					AddRow("1"),
			)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.RegistrationRequest{
			Username: "gam6itko@gmail.com",
			Password: "correct password",
		}
		resp, err := svc.Registration(context.TODO(), &req)
		require.Nil(t, resp)
		require.Error(t, err)
		require.EqualError(t, err, "rpc error: code = AlreadyExists desc = user already exists")
	})

	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.
			ExpectQuery("SELECT (.+) FROM `user`").
			WillReturnRows(
				sqlmock.NewRows([]string{"cnt"}).
					AddRow("0"),
			)
		// Сохраняем пользователя в БД.
		mock.
			ExpectExec("INSERT INTO `user`").
			WithArgs("gam6itko@gmail.com", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(2, 1))

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.RegistrationRequest{
			Username: "gam6itko@gmail.com",
			Password: "correct password",
		}
		resp, err := svc.Registration(context.TODO(), &req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

func TestAuthServerImpl_Login(t *testing.T) {
	t.Run("username is too short", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.LoginRequest{}
		resp, err := svc.Login(context.TODO(), &req)
		require.Nil(t, resp)
		require.Error(t, err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = username is too short")
	})

	t.Run("password is too short", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.LoginRequest{
			Username: "gam6itko@gmail.com",
		}
		resp, err := svc.Login(context.TODO(), &req)
		require.Nil(t, resp)
		require.Error(t, err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = password is too short")
	})

	t.Run("wrong password", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.
			ExpectQuery("SELECT (.+) FROM `user`").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "parsing_url"}).
					AddRow(
						"1",
						[]byte("wrong password"),
					),
			)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.LoginRequest{
			Username: "gam6itko@gmail.com",
			Password: "correct password",
		}
		resp, err := svc.Login(context.TODO(), &req)
		require.Nil(t, resp)
		require.Error(t, err)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = username or password is incorrect")
	})

	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		// Пароль в БД.
		password, err := bcrypt.GenerateFromPassword([]byte("correct password"), bcrypt.MinCost)
		require.NoError(t, err)

		mock.
			ExpectQuery("SELECT (.+) FROM `user`").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "parsing_url"}).
					AddRow("1", password),
			)

		issuer := jwt.NewIssuer()
		svc := NewAuthServerImpl(db, issuer)
		req := proto.LoginRequest{
			Username: "gam6itko@gmail.com",
			Password: "correct password",
		}
		resp, err := svc.Login(context.TODO(), &req)
		require.NotNil(t, resp)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Token)

		// Проверяем токен на корректность.
		claims, err := issuer.Parse(resp.Token)
		require.NoError(t, err)
		assert.NotNil(t, claims.RegisteredClaims.ExpiresAt)
		assert.NotNil(t, claims.RegisteredClaims.IssuedAt)
		assert.Equal(t, "server", claims.RegisteredClaims.Issuer)
		assert.Equal(t, uint64(1), claims.UserID)
	})
}
