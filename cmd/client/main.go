package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"

	"github.com/gam6itko/goph-keeper/internal/client/config"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	grpcServer "github.com/gam6itko/goph-keeper/internal/client/server/grpc"
	"github.com/gam6itko/goph-keeper/internal/client/tui/root"
	"github.com/gam6itko/goph-keeper/proto"
)

func main() {
	cfg := config.Load()
	creds, err := credentials.NewClientTLSFromFile(cfg.GRPC.TLS.CertPEM, cfg.GRPC.TLS.ServerHost)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	// Устанавливаем соединение с сервером.
	conn, err := grpc.NewClient(
		cfg.GRPC.ServerAddr,
		grpc.WithTransportCredentials(creds),
	)

	s := grpcServer.New(
		proto.NewAuthClient(conn),
		proto.NewKeeperClient(conn),
	)

	memStorage := masterkey.NewMemGuardStorage()
	crypt := encrypt.NewAESCrypt()
	p := tea.NewProgram(
		root.New(s, memStorage, crypt),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
