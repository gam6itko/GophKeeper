package errmsg

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type (
//	GotoModelMsg struct {
//		Model tea.Model
//	}
)

type Model struct {
	err       error
	returnCmd tea.Cmd
}

func New(err error, returnCmd tea.Cmd) *Model {
	return &Model{
		err:       err,
		returnCmd: returnCmd,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, m.returnCmd
	}

	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("Error. %s\nPress any key to continue.", m.err.Error())
}
