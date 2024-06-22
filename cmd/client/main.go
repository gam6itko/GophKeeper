package main

import (
	"fmt"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/root"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	s := &server.MockServer{}
	p := tea.NewProgram(root.New(s), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
