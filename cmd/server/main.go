// Сервер должен реализовывать следующую бизнес-логику:
// - регистрация, аутентификация и авторизация пользователей;
// - хранение приватных данных;
// - синхронизация данных между несколькими авторизованными клиентами одного владельца;
// - передача приватных данных владельцу по запросу.

package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	jwt_inter "github.com/gam6itko/goph-keeper/internal/server/interceptors/jwt"
	"github.com/gam6itko/goph-keeper/internal/server/jwt"
	"github.com/gam6itko/goph-keeper/internal/server/service"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"time"

	"github.com/gam6itko/goph-keeper/internal/server/config"
	"github.com/gam6itko/goph-keeper/proto"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	initDb(db)

	// определяем порт для сервера
	listen, err := net.Listen("tcp", cfg.GRPC.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile(cfg.GRPC.TLS.CertPEM, cfg.GRPC.TLS.KeyPEM)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	issuer := jwt.NewIssuer(
		jwt.WithKey([]byte(cfg.JWT.Secret)),
		jwt.WithExpiresIn(time.Duration(cfg.JWT.ExpiresInSeconds)),
	)
	// Создаём gRPC-сервер.
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(jwt_inter.New(issuer).Intercept),
	)
	// Регистрируем сервисы.
	proto.RegisterAuthServer(s, service.NewAuthServerImpl(db, issuer))
	proto.RegisterKeeperServer(s, service.NewKeeperImpl(db))

	fmt.Println("gRPC server listening on " + cfg.GRPC.ServerAddr)
	// получаем запрос gRPC
	if err = s.Serve(listen); err != nil {
		log.Printf("server stop. err: %s", err)
	}
}
