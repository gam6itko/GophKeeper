// Сервер должен реализовывать следующую бизнес-логику:
// - регистрация, аутентификация и авторизация пользователей;
// - хранение приватных данных;
// - синхронизация данных между несколькими авторизованными клиентами одного владельца;
// - передача приватных данных владельцу по запросу.

package main

import (
	"database/sql"
	"fmt"
	"github.com/gam6itko/goph-keeper/internal/server"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"

	"github.com/gam6itko/goph-keeper/internal/server/config"
	"github.com/gam6itko/goph-keeper/internal/server/interceptors/jwt"
	"github.com/gam6itko/goph-keeper/proto"
)

func main() {
	cfg := config.Load()

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

	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	// Создаём gRPC-сервер без зарегистрированной службы.
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(jwt.New().Intercept),
	)
	// Регистрируем сервисы.
	proto.RegisterAuthServer(s, server.NewAuthServerImpl(db))
	proto.RegisterKeeperServer(s, server.NewKeeperImpl(db))

	fmt.Println("gRPC server listening on " + cfg.GRPC.ServerAddr)
	// получаем запрос gRPC
	if err = s.Serve(listen); err != nil {
		log.Printf("server stop. err: %s", err)
	}
}

//todo
//	- Схема БД, создание БД при запуске приложения если БД не инициализирована.
//  - Запуск gRPC сервера.
//	- Ендпоинты для Регистрации, Аутентифакации, Сохранения данных

//todo Login - jwt token
