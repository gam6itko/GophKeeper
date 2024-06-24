package main

import (
	"fmt"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/root"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	s := &server.MockServer{}
	//memStorage := masterkey.NewMemGuardStorage()
	memStorage := &masterkey.SimpleStorage{}
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
