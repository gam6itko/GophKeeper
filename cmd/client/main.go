package main

import (
	"fmt"
	"github.com/gam6itko/goph-keeper/internal/client/tui/root"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(root.NewScreen(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
